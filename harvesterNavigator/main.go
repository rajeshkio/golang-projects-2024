package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeconfig(kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}
	return clientset, nil
}

func main() {
	kubeconfigPath := flag.String("kubeconfig", "", "Path to kubeconfig file (optional, falls back to $KUBECONFIG)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: ./harvester_vm_info [--kubeconfig <path>] <vm-name>")
		os.Exit(1)
	}

	clientset, err := getKubeconfig(*kubeconfigPath)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	vmName := flag.Arg(0)
	namespace := "default"

	// fmt.Printf("Fetching details of the VM: %s\n", vmName)
	vm, err := clientset.RESTClient().Get().
		AbsPath("apis/kubevirt.io/v1").
		Namespace(namespace).
		Resource("virtualmachines").
		Name(vmName).
		Do(context.TODO()).Raw()
	if err != nil {
		fmt.Printf("Couldn't fetch the VM: %v", err)
	}

	// fmt.Printf("\nVM Name: %s\n", vm)
	var vmData map[string]interface{}
	if err := json.Unmarshal(vm, &vmData); err != nil {
		log.Fatalf("Error parsing VM JSON: %v", err)
	}

	//var prettyJSON bytes.Buffer

	//	if err := json.Indent(&prettyJSON, vm, "", " "); err != nil {
	//		fmt.Printf("failed to format JSON: %v", err)
	//	}
	//	fmt.Println(prettyJSON.String())
	metadata := vmData["metadata"].(map[string]interface{})
	//fmt.Println("VM Name:", name)
	//nameTest := metadata["name"].(string)
	//fmt.Println("VM name: ", nameTest)
	//fmt.Println("")

	annotations := metadata["annotations"].(map[string]interface{})
	volumeClainTemplateInterface := annotations["harvesterhci.io/volumeClaimTemplates"].(string)
	//fmt.Println("volumeClaimTemplateInterface: ", volumeClainTemplateInterface)
	//fmt.Println("")

	var volumeClaimTemplate []map[string]interface{}
	json.Unmarshal([]byte(volumeClainTemplateInterface), &volumeClaimTemplate)

	//fmt.Println("volumeClaimTemplate: ", volumeClaimTemplate)
	//fmt.Println("")

	vmImageId := volumeClaimTemplate[0]["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["harvesterhci.io/imageId"]
	vmStorageClass := volumeClaimTemplate[0]["spec"].(map[string]interface{})["storageClassName"]

	vmVolumes := vmData["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["volumes"].([]interface{})

	var claimNames []string

	for _, vol := range vmVolumes {
		volume := vol.(map[string]interface{})
		if pvc, exists := volume["persistentVolumeClaim"]; exists {
			persistentVolumeClaim := pvc.(map[string]interface{})
			claimName := persistentVolumeClaim["claimName"].(string)
			claimNames = append(claimNames, claimName)
		}
	}

	raw, err := clientset.RESTClient().
		Get().
		AbsPath("/apis").
		Do(context.TODO()).
		Raw()

	if err != nil {
		log.Fatalf("Error fetching API groups: %v", err)
	}

	// Format JSON output
	var formattedJSON map[string]interface{}
	if err := json.Unmarshal(raw, &formattedJSON); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	//prettyJSON, _ := json.MarshalIndent(formattedJSON, "", "  ")
	//fmt.Println(string(prettyJSON))

	//fmt.Printf("Fetching details of PVC: %s\n", claimNames[0])
	pvcDetails, err := clientset.RESTClient().Get().
		AbsPath("api/v1").
		Namespace(namespace).
		Resource("persistentvolumeclaims").
		Name(claimNames[0]).
		Do(context.TODO()).Raw()
	if err != nil {
		fmt.Printf("Couldn't fetch the volume: %v", err)
	}

	var pvcData map[string]interface{}
	if err := json.Unmarshal(pvcDetails, &pvcData); err != nil {
		log.Fatalf("Error parsing volumeName JSON: %v", err)
	}
	pvcSpec := pvcData["spec"].(map[string]interface{})
	//fmt.Printf("PVC Spec: %+v\n", pvcSpec)

	volumeName := pvcSpec["volumeName"].(string)
	//fmt.Printf("Volume Name: %+v\n", volumeName)

	//fmt.Printf("Fetching details of volumes: %s\n", volumeName)
	volumeDetails, err := clientset.RESTClient().Get().
		AbsPath("apis/longhorn.io/v1beta2").
		Namespace("longhorn-system").
		Resource("volumes").
		Name(volumeName).
		Do(context.TODO()).Raw()
	if err != nil {
		fmt.Printf("Couldn't fetch the volume: %v", err)
	}

	var volumeData map[string]interface{}
	if err := json.Unmarshal(volumeDetails, &volumeData); err != nil {
		log.Fatalf("Error parsing volumeName JSON: %v", err)
	}

	//fmt.Printf(" volume: %v", volumeData)

	volumeStatus := volumeData["status"].(map[string]interface{})
	workloadStatus := volumeStatus["kubernetesStatus"].(map[string]interface{})["workloadsStatus"].([]interface{})
	var podName string
	if len(workloadStatus) > 0 {
		firstWorkload := workloadStatus[0].(map[string]interface{})

		// Now you can get the podName
		podName = firstWorkload["podName"].(string)
		//fmt.Printf("Pod Name: %s\n", podName)
	}

	//fmt.Printf("Fetching details of volumeAttachment: %s\n", volumeName)
	volumeAttachmentDetails, err := clientset.RESTClient().Get().
		AbsPath("apis/longhorn.io/v1beta2").
		Namespace("longhorn-system").
		Resource("volumeattachments").
		Name(volumeName).
		Do(context.TODO()).Raw()
	if err != nil {
		fmt.Printf("Couldn't fetch the volume: %v", err)
	}

	var volumeAttachmentData map[string]interface{}
	if err := json.Unmarshal(volumeAttachmentDetails, &volumeAttachmentData); err != nil {
		log.Fatalf("Error parsing volumeAttachmentDetails JSON: %v", err)
	}

	//fmt.Printf(" volume: %v", volumeData)

	volumeAttachmentTickets := volumeAttachmentData["spec"].(map[string]interface{})["attachmentTickets"].(map[string]interface{})

	fmt.Println("VM Name:", vmName)
	fmt.Println("VM Image ID:", vmImageId)
	fmt.Printf("Pod Name: %s\n", podName)
	fmt.Println("VM Storage Class:", vmStorageClass)
	fmt.Println("PVC Claim Names:", claimNames)
	fmt.Println("Volume Name:", volumeName)
	fmt.Println(volumeAttachmentTickets)
	//["workloadsStatus"].(map[string]interface{})["podName"])

}
