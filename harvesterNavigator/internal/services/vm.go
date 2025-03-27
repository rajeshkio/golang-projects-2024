package vm

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

func FetchVMData(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	vm, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get vmdata: %s", err)
	}

	var vmData map[string]interface{}
	err = json.Unmarshal(vm, &vmData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall vmdata: %s", err)
	}
	return vmData, nil
}

func ParseVMMetaData(vmData map[string]interface{}, vmInfo *types.VMInfo) error {
	metadata, ok := vmData["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("metadata field missing or not an object")
	}
	//	fmt.Println(metadata["name"])
	//fmt.Println()

	annotations, ok := metadata["annotations"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("annotations field missing or not an object")
	}
	volumeClaimTemplateStr, ok := annotations["harvesterhci.io/volumeClaimTemplates"].(string)
	if !ok {
		return fmt.Errorf("volumeClaimTemplates annotation missing or not a string")
	}
	var volumeClaimTemplates []map[string]interface{}
	err := json.Unmarshal([]byte(volumeClaimTemplateStr), &volumeClaimTemplates)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the VM volumeclaim template: %w", err)
	}

	if len(volumeClaimTemplates) == 0 {
		return fmt.Errorf("no volume claim templates found")
	}

	template := volumeClaimTemplates[0]

	templateMetadata, ok := template["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("template metadata missing or not an object")
	}

	templateMetaAnnotation, ok := templateMetadata["annotations"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("template annotations missing or not an object")
	}
	imageID, ok := templateMetaAnnotation["harvesterhci.io/imageId"].(string)
	if !ok {
		return fmt.Errorf("imageId missing or not a string")
	}

	vmInfo.ImageId = imageID

	templateSpec, ok := template["spec"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("template spec missing or not an object")
	}

	storageClass, ok := templateSpec["storageClassName"].(string)
	if !ok {
		return fmt.Errorf("storageClassName missing or not a string")
	}

	vmInfo.StorageClass = storageClass

	return nil
}

func DisplayVMInfo(info *types.VMInfo) {
	fmt.Println("VM Name:", info.Name)
	fmt.Println("VM Image ID:", info.ImageId)
	fmt.Println("Pod Name:", info.PodName)
	fmt.Println("VM Storage Class:", info.StorageClass)
	fmt.Println("PVC Claim Names:", info.ClaimNames)
	fmt.Println("Volume Name:", info.VolumeName)
	fmt.Println(info.AttachmentInfo)

}
