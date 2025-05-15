package pod

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

func FetchPodDetails(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	vm, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get podDetails: %s", err)
	}

	var podData map[string]interface{}
	err = json.Unmarshal(vm, &podData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall poddetails: %s", err)
	}
	return podData, nil

}

func ParsePodData(podData map[string]interface{}) ([]types.PodInfo, error) {
	var podInfos []types.PodInfo
	podMetadata := podData["metadata"].(map[string]interface{})
	podSpec, _ := podData["spec"].(map[string]interface{})
	podStatus, _ := podData["status"].(map[string]interface{})
	ownerRefName := ""
	if ownerRefs, ok := podMetadata["ownerReferences"].([]interface{}); ok && ownerRefs != nil && len(ownerRefs) > 0 {
		if ownerRef, ok := ownerRefs[0].(map[string]interface{}); ok && ownerRef != nil {
			if name, ok := ownerRef["name"].(string); ok {
				ownerRefName = name
			}
		}
	}

	nodeName, ok := podSpec["nodeName"].(string)
	if !ok {
		fmt.Println("cannot get nodename")
	}
	status, ok := podStatus["phase"].(string)
	if !ok {
		fmt.Println("cannot get phase")
	}

	podInfo := types.PodInfo{
		VMI:    ownerRefName,
		NodeID: nodeName,
		Status: status,
	}
	podInfos = append(podInfos, podInfo)
	return podInfos, nil
}
