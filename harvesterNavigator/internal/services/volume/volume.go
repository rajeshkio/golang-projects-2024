package volume

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func FetchVolumeDetails(client *kubernetes.Clientset, name, absPath, namespace, resource string) (map[string]interface{}, error) {
	volume, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Name(name).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %s", err)
	}

	var volumeData map[string]interface{}
	err = json.Unmarshal(volume, &volumeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall volumeDetails: %s", err)
	}
	return volumeData, nil
}

func ParseVolumeSpec(volumeData map[string]interface{}) (string, error) {
	pvcSpec := volumeData["spec"].(map[string]interface{})
	nodeID := pvcSpec["nodeID"].(string)

	fmt.Println("nodeID")

	return nodeID, nil
}

func GetPodFromVolume(volumeData map[string]interface{}) (string, error) {
	volumeStatus := volumeData["status"].(map[string]interface{})
	workloadStatus := volumeStatus["kubernetesStatus"].(map[string]interface{})["workloadsStatus"].([]interface{})
	var podName string
	if len(workloadStatus) > 0 {
		firstWorkload := workloadStatus[0].(map[string]interface{})

		// Now you can get the podName
		podName = firstWorkload["podName"].(string)
		//fmt.Printf("Pod Name: %s\n", podName)
	}
	return podName, nil

}
