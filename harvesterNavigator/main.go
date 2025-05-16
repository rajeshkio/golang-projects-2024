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
	"github.com/rk280392/harvesterNavigator/internal/services/pod"
	"github.com/rk280392/harvesterNavigator/internal/services/replicas"
	vm "github.com/rk280392/harvesterNavigator/internal/services/vm"
	vmi "github.com/rk280392/harvesterNavigator/internal/services/vmi"
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

func logNotFound(resourceType, name, namespace string, err error) {
	log.Printf(
		"\nError: %s %q not found in namespace %q.\nCheck if the %s exists and that the namespace is correct.\nDetails: %v",
		resourceType, name, namespace, resourceType, err,
	)
}

// Set missing resource and display information before exiting
func handleResourceError(resourceType string, vmInfo *types.VMInfo) {
	vmInfo.MissingResource = resourceType
	display.DisplayVMInfo(vmInfo)
	os.Exit(1)
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

	// Initialize the VM info structure
	vmInfo := &types.VMInfo{Name: vmName}

	// Fetch VM data
	absPath := "apis/kubevirt.io/v1"
	resource := "virtualmachines"
	vmData, err := vm.FetchVMData(clientset, vmName, absPath, namespace, resource)
	if err != nil {
		logNotFound("VM", vmName, namespace, err)
		handleResourceError("VM", vmInfo)
	}

	// Parse VM metadata
	err = vm.ParseVMMetaData(vmData, vmInfo)
	if err != nil {
		log.Printf("Failed to parse VM metadata: %s", err)
		handleResourceError("VM METADATA", vmInfo)
	}

	// Fetch PVC data
	pvcAPIPath := "/api/v1"
	pvcResource := "persistentvolumeclaims"
	pvcData, err := pvc.FetchPVCData(clientset, vmInfo.ClaimNames, pvcAPIPath, namespace, pvcResource)
	if err != nil {
		logNotFound("PVC", vmInfo.ClaimNames, namespace, err)
		handleResourceError("PVC", vmInfo)
	}

	// Parse PVC spec
	volumeName, err := pvc.ParsePVCSpec(pvcData)
	if err != nil {
		log.Printf("Failed to parse PVC spec: %s", err)
		handleResourceError("PVC SPEC", vmInfo)
	}
	vmInfo.VolumeName = volumeName

	// Parse PVC status
	status, err := pvc.ParsePVCStatus(pvcData)
	if err != nil {
		log.Printf("Failed to parse PVC status: %s", err)
		handleResourceError("PVC STATUS", vmInfo)
	}
	vmInfo.PVCStatus = status

	// Fetch volume details
	volumeAPIPath := "apis/longhorn.io/v1beta2"
	volNamespace := "longhorn-system"
	volumeResource := "volumes"
	volumeDetails, err := volume.FetchVolumeDetails(clientset, volumeName, volumeAPIPath, volNamespace, volumeResource)
	if err != nil {
		logNotFound("Volume", volumeName, volNamespace, err)
		handleResourceError("VOLUME", vmInfo)
	}

	// Get pod name from volume
	podName, err := volume.GetPodFromVolume(volumeDetails)
	if err != nil {
		log.Printf("Failed to get pod name from volume status: %s", err)
		handleResourceError("POD NAME", vmInfo)
	}
	vmInfo.PodName = podName

	// Fetch pod details
	podApiPath := "/api/v1"
	podResource := "pods"
	podData, err := pod.FetchPodDetails(clientset, podName, podApiPath, namespace, podResource)
	if err != nil {
		logNotFound("POD", podName, namespace, err)
		handleResourceError("POD", vmInfo)
	}

	// Check if podData is nil to prevent panic
	if podData == nil {
		log.Printf("Error: Pod data is unexpectedly nil")
		handleResourceError("POD DATA", vmInfo)
	}

	// Parse pod data with panic protection
	var ownerRef []types.PodInfo
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic during pod data parsing: %v", r)
				handleResourceError("POD DATA", vmInfo)
			}
		}()

		var err error
		ownerRef, err = pod.ParsePodData(podData)
		if err != nil {
			log.Printf("Failed to parse pod data: %s", err)
			handleResourceError("POD DATA", vmInfo)
		}
	}()

	vmInfo.PodInfo = ownerRef

	// Extract VMI name
	vmiName := ""
	if len(vmInfo.PodInfo) > 0 {
		for _, pod := range vmInfo.PodInfo {
			if pod.VMI != "" {
				vmiName = pod.VMI
				break
			}
		}
	}

	if vmiName == "" {
		log.Printf("Error: No VMI name found in pod data")
		handleResourceError("VMI NAME", vmInfo)
	}

	// Fetch VMI details
	vmiApiPath := "apis/kubevirt.io/v1"
	vmiResource := "virtualmachineinstances"
	vmiData, err := vmi.FetchVMIDetails(clientset, vmiName, vmiApiPath, namespace, vmiResource)
	if err != nil {
		logNotFound("VMI", vmiName, namespace, err)
		handleResourceError("VMI", vmInfo)
	}

	// Parse VMI data
	vmiStatus, err := vmi.ParseVMIData(vmiData)
	if err != nil {
		log.Printf("Failed to parse VMI data: %s", err)
		handleResourceError("VMI DATA", vmInfo)
	}
	vmInfo.VMIInfo = vmiStatus

	// Find replica details
	replicaAPIPath := "apis/longhorn.io/v1beta2"
	replicaNamespace := "longhorn-system"
	replicaResource := "replicas"
	relatedReplicas, err := replicas.FindReplicaDetails(clientset, volumeName, replicaAPIPath, replicaNamespace, replicaResource)
	if err != nil {
		log.Printf("Failed to get replica details: %s", err)
		handleResourceError("REPLICAS", vmInfo)
	}
	vmInfo.ReplicaInfo = relatedReplicas

	// Find engine details - this is optional
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

	// Display all collected information
	display.DisplayVMInfo(vmInfo)
}
