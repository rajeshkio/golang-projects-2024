package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	kubeclient "github.com/rk280392/harvesterNavigator/internal/client"
	types "github.com/rk280392/harvesterNavigator/internal/models"
	"github.com/rk280392/harvesterNavigator/internal/services/engine"
	pvc "github.com/rk280392/harvesterNavigator/internal/services/longhornPVC"
	"github.com/rk280392/harvesterNavigator/internal/services/replicas"
	vm "github.com/rk280392/harvesterNavigator/internal/services/vm"
	volume "github.com/rk280392/harvesterNavigator/internal/services/volume"
	display "github.com/rk280392/harvesterNavigator/pkg/display"
)

func main() {
	kubeconfigPath := flag.String("kubeconfig", "", "Path to kubeconfig file (optional, falls back to $KUBECONFIG)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: ./harvester_vm_info [--kubeconfig <path>] <vm-name> [--namespace <name>]")
		os.Exit(1)
	}

	clientset, err := kubeclient.NewClient(*kubeconfigPath)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	vmName := flag.Arg(0)
	namespace := flag.Arg(1)
	absPath := "apis/kubevirt.io/v1"
	resource := "virtualmachines"

	vmData, err := vm.FetchVMData(clientset, vmName, absPath, namespace, resource)
	if err != nil {
		log.Fatalf("failed to fetch the VM Data: %s", err)
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
		log.Fatalf("failed to fetch the VM Data: %s", err)
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
	volumeResourceName := volumeName
	volumeDetails, err := volume.FetchVolumeDetails(clientset, volumeResourceName, volumeAPIPath, volNamespace, volumeResource)
	if err != nil {
		log.Fatalf("failed to get volumeDetails: %s", err)
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

	relatedReplicas, err := replicas.FindReplicaDetails(clientset, volumeResourceName, replicaAPIPath, replicaNamespace, replicaResource)
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
