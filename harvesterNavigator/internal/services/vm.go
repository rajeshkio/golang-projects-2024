package vm

import (
	"context"
	"encoding/json"
	"fmt"

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
