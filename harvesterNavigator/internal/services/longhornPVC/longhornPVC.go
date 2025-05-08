package longhornPVC

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func FetchPVCData(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	vm, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get pvcdata: %s", err)
	}

	var pvcData map[string]interface{}
	err = json.Unmarshal(vm, &pvcData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall pvcdata: %s", err)
	}
	return pvcData, nil
}

func ParsePVCSpec(pvcData map[string]interface{}) (string, error) {
	pvcSpec := pvcData["spec"].(map[string]interface{})
	volumeName := pvcSpec["volumeName"].(string)

	return volumeName, nil
}

func ParsePVCStatus(pvcData map[string]interface{}) (string, error) {
	pvcStatus := pvcData["status"].(map[string]interface{})
	status := pvcStatus["phase"].(string)

	return status, nil
}
