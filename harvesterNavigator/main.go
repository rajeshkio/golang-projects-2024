package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	kubeclient "github.com/rk280392/harvesterNavigator/internal/client"
	types "github.com/rk280392/harvesterNavigator/internal/models"
	vm "github.com/rk280392/harvesterNavigator/internal/services"
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

	//fmt.Println(vmData)

	vmInfo := &types.VMInfo{Name: vmName}
	err = vm.ParseVMMetaData(vmData, vmInfo)

	vm.DisplayVMInfo(vmInfo)

}
