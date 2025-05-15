package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	kubeclient "github.com/rk280392/harvesterNavigator/internal/client"
	types "github.com/rk280392/harvesterNavigator/internal/models"
	"github.com/rk280392/harvesterNavigator/internal/services/engine"
	pvc "github.com/rk280392/harvesterNavigator/internal/services/longhornPVC"
	"github.com/rk280392/harvesterNavigator/internal/services/replicas"
	vm "github.com/rk280392/harvesterNavigator/internal/services/vm"
	volume "github.com/rk280392/harvesterNavigator/internal/services/volume"
	display "github.com/rk280392/harvesterNavigator/pkg/display"
	flag "github.com/spf13/pflag"
)

func defaultKubeconfigPath() string {
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env
	}
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, ".kube", "config")
}
func getNamespace(cliNamespace string) string {
	if cliNamespace != "" {
		return cliNamespace
	}
	if env := os.Getenv("NAMESPACE"); env != "" {
		return env
	}
	return "default"
}

func fatalNotFound(resourceType, name, namespace string, err error) {
	log.Fatalf(
		"\nError: %s %q not found in namespace %q.\nCheck if the %s exists and that the namespace is correct.\nDetails: %v",
		resourceType, name, namespace, resourceType, err,
	)
}

func main() {

	// Define optional flags
	kubeconfig := flag.StringP("kubeconfig", "k", defaultKubeconfigPath(), "Path to kubeconfig file (optional)")
	cliNamespace := flag.StringP("namespace", "n", "", "Namespace of the VM (optional, or export NAMESPACE env var)")

	// Override default usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <vm-name>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Validate positional arg: VM name
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: VM name is required.")
		flag.Usage()
		os.Exit(1)
	}

	vmName := flag.Arg(0)
	namespace := getNamespace(*cliNamespace)

	if _, err := os.Stat(*kubeconfig); os.IsNotExist(err) {
		log.Fatalf("Error: kubeconfig file not found at '%s'", *kubeconfig)
	}

	clientset, err := kubeclient.NewClient(*kubeconfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	absPath := "apis/kubevirt.io/v1"
	resource := "virtualmachines"

	vmData, err := vm.FetchVMData(clientset, vmName, absPath, namespace, resource)
	if err != nil {
		fatalNotFound("VM", vmName, namespace, err)
	}

	vmInfo := &types.VMInfo{Name: vmName}
	err = vm.ParseVMMetaData(vmData, vmInfo)
	if err != nil {
		log.Fatalf("failed to parse VMMetadata: %s", err)
	}

	pvcAPIPath := "/api/v1"
	pvcResource := "persistentvolumeclaims"
	pvcData, err := pvc.FetchPVCData(clientset, vmInfo.ClaimNames, pvcAPIPath, namespace, pvcResource)
	if err != nil {
		fatalNotFound("PVC", vmInfo.ClaimNames, namespace, err)
	}

	volumeName, err := pvc.ParsePVCSpec(pvcData)
	if err != nil {
		log.Fatalf("failed to parse PVC Spec: %s", err)
	}
	vmInfo.VolumeName = volumeName

	status, err := pvc.ParsePVCStatus(pvcData)
	if err != nil {
		log.Fatalf("failed to parse PVC Spec: %s", err)
	}
	vmInfo.PVCStatus = status

	volumeAPIPath := "apis/longhorn.io/v1beta2"
	volNamespace := "longhorn-system"
	volumeResource := "volumes"
	volumeDetails, err := volume.FetchVolumeDetails(clientset, volumeName, volumeAPIPath, volNamespace, volumeResource)
	if err != nil {
		fatalNotFound("Volume", volumeName, volNamespace, err)
	}

	//volumeInfo := &types.VolumeInfo{Name: volumeName}
	podName, err := volume.GetPodFromVolume(volumeDetails)
	if err != nil {
		log.Fatalf("failed to get podname from volume Status: %s", err)
	}
	vmInfo.PodName = podName

	replicaAPIPath := "apis/longhorn.io/v1beta2"
	replicaNamespace := "longhorn-system"
	replicaResource := "replicas"

	relatedReplicas, err := replicas.FindReplicaDetails(clientset, volumeName, replicaAPIPath, replicaNamespace, replicaResource)
	if err != nil {
		log.Fatalf("failed to get replica from volume: %s", err)
	}
	vmInfo.ReplicaInfo = relatedReplicas

	engineAPIPath := "apis/longhorn.io/v1beta2"
	engineNamespace := "longhorn-system"
	engineResource := "engines"

	engineInfos, err := engine.FindEngineDetails(clientset, vmInfo.VolumeName, engineAPIPath, engineNamespace, engineResource)
	if err != nil {
		// Log the error but continue - engine info is optional
		log.Printf("Warning: failed to get engine details: %s", err)
	} else {
		vmInfo.EngineInfo = engineInfos
	}

	display.DisplayVMInfo(vmInfo)

}
