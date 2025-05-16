package vm

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

// FetchVMData retrieves virtual machine data from the Kubernetes API.
func FetchVMData(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	vm, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get VM data: %w", err)
	}

	var vmData map[string]interface{}
	if err := json.Unmarshal(vm, &vmData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VM data: %w", err)
	}

	return vmData, nil
}

// ParseVMMetaData extracts relevant information from VM data and populates the VMInfo struct.
func ParseVMMetaData(vmData map[string]interface{}, vmInfo *types.VMInfo) error {
	// Try to extract metadata
	metadata, ok := vmData["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("metadata field missing or not an object")
	}

	// Extract status information
	statusRaw, ok := vmData["status"]
	if ok {
		status, ok := statusRaw.(map[string]interface{})
		if ok {
			// Extract printable status if available
			if printableStatus, ok := status["printableStatus"]; ok && printableStatus != nil {
				if printableStatusStr, ok := printableStatus.(string); ok {
					vmInfo.PrintableStatus = printableStatusStr
				}
			}

			// Try to parse conditions if available
			if conditionRaw, ok := status["conditions"]; ok && conditionRaw != nil {
				if conditionsArray, ok := conditionRaw.([]interface{}); ok {
					for _, conditionRaw := range conditionsArray {
						condition, ok := conditionRaw.(map[string]interface{})
						if !ok {
							continue // Skip invalid entries
						}

						if typeVal, ok := condition["reason"]; ok && typeVal != nil {
							if typeStr, ok := typeVal.(string); ok {
								vmInfo.VMStatusReason = typeStr
							}
						}

						if statusVal, ok := condition["status"]; ok && statusVal != nil {
							if statusStr, ok := statusVal.(string); ok {
								vmInfo.VMStatus = statusStr
							}
						}
					}
				}
			}
		}
	}

	// Try to extract annotations
	annotationsRaw, ok := metadata["annotations"]
	if !ok {
		return fmt.Errorf("annotations field missing")
	}

	annotations, ok := annotationsRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("annotations is not an object")
	}

	// Try to extract volume claim templates
	volumeClaimTemplateRaw, ok := annotations["harvesterhci.io/volumeClaimTemplates"]
	if !ok {
		return fmt.Errorf("volumeClaimTemplates annotation missing")
	}

	volumeClaimTemplateStr, ok := volumeClaimTemplateRaw.(string)
	if !ok {
		return fmt.Errorf("volumeClaimTemplates is not a string")
	}

	var volumeClaimTemplates []map[string]interface{}
	if err := json.Unmarshal([]byte(volumeClaimTemplateStr), &volumeClaimTemplates); err != nil {
		return fmt.Errorf("failed to unmarshal volume claim templates: %w", err)
	}

	if len(volumeClaimTemplates) == 0 {
		return fmt.Errorf("no volume claim templates found")
	}

	// Process the first template
	template := volumeClaimTemplates[0]

	templateMetadataRaw, ok := template["metadata"]
	if !ok {
		return fmt.Errorf("template metadata missing")
	}

	templateMetadata, ok := templateMetadataRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("template metadata is not an object")
	}

	// Extract PVC claim name
	if nameRaw, ok := templateMetadata["name"]; ok && nameRaw != nil {
		if name, ok := nameRaw.(string); ok {
			vmInfo.ClaimNames = name
		} else {
			return fmt.Errorf("PVC claim name is not a string")
		}
	} else {
		return fmt.Errorf("PVC claim name missing")
	}

	// Try to extract template metadata annotations
	if templateMetaAnnotationRaw, ok := templateMetadata["annotations"]; ok && templateMetaAnnotationRaw != nil {
		if templateMetaAnnotation, ok := templateMetaAnnotationRaw.(map[string]interface{}); ok {
			if imageIDRaw, ok := templateMetaAnnotation["harvesterhci.io/imageId"]; ok && imageIDRaw != nil {
				if imageID, ok := imageIDRaw.(string); ok {
					vmInfo.ImageId = imageID
				}
			}
		}
	}

	// Try to extract template spec
	templateSpecRaw, ok := template["spec"]
	if !ok {
		return fmt.Errorf("template spec missing")
	}

	templateSpec, ok := templateSpecRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("template spec is not an object")
	}

	// Try to extract storage class
	if storageClassRaw, ok := templateSpec["storageClassName"]; ok && storageClassRaw != nil {
		if storageClass, ok := storageClassRaw.(string); ok {
			vmInfo.StorageClass = storageClass
		} else {
			return fmt.Errorf("storageClassName is not a string")
		}
	} else {
		return fmt.Errorf("storageClassName missing")
	}

	return nil
}
