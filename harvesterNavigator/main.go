package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	kubeclient "github.com/rk280392/harvesterNavigator/internal/client"
	types "github.com/rk280392/harvesterNavigator/internal/models"
	pvc "github.com/rk280392/harvesterNavigator/internal/services/longhornPVC"
	vm "github.com/rk280392/harvesterNavigator/internal/services/vm"
	volume "github.com/rk280392/harvesterNavigator/internal/services/volume"
)

func main() {
	kubeconfigPath := flag.String("kubeconfig", "", "Path to kubeconfig file (optional, falls back to $KUBECONFIG)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: ./harvester_vm_info [--kubeconfig <path>] <vm-name>")
		os.Exit(1)
	}

	clientset, err := kubeclient.NewClient(*kubeconfigPath)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	vmName := flag.Arg(0)
	namespace := "default"
	absPath := "apis/kubevirt.io/v1"
	resource := "virtualmachines"

	vmData, err := vm.FetchVMData(clientset, vmName, absPath, namespace, resource)
	if err != nil {
		log.Fatalf("failed to fetch the VM Data: %s", err)
	}

	vmInfo := &types.VMInfo{Name: vmName}
	err = vm.ParseVMMetaData(vmData, vmInfo)

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

	//volumeInfo := &types.VolumeInfo{Name: volumeName}
	podName, err := volume.GetPodFromVolume(volumeDetails)
	if err != nil {
		log.Fatalf("failed to get podname from volume Status: %s", err)
	}
	vmInfo.PodName = podName
	vm.DisplayVMInfo(vmInfo)
}
