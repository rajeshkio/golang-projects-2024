package replicas

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

func FindReplicaDetails(client *kubernetes.Clientset, volumeName, absPath, namespace, resource string) ([]types.ReplicaInfo, error) {
	replicas, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get replicas: %s", err)
	}

	var replicaData map[string]interface{}
	err = json.Unmarshal(replicas, &replicaData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall replicaData: %s", err)
	}
	items, ok := replicaData["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get items from replicas response")
	}

	var relatedReplicas []map[string]interface{}
	for _, item := range items {
		replica, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		spec, ok := replica["spec"].(map[string]interface{})
		if !ok {
			continue
		}
		if replicaVolumename, ok := spec["volumeName"].(string); ok && replicaVolumename == volumeName {
			relatedReplicas = append(relatedReplicas, replica)
		}
	}
	var replicaInfos []types.ReplicaInfo
	for _, replica := range relatedReplicas {
		metadata, _ := replica["metadata"].(map[string]interface{})
		spec, _ := replica["spec"].(map[string]interface{})
		status, _ := replica["status"].(map[string]interface{})

		ownerRefName := ""
		if ownerRefs, ok := metadata["ownerReferences"].([]interface{}); ok && ownerRefs != nil && len(ownerRefs) > 0 {
			if ownerRef, ok := ownerRefs[0].(map[string]interface{}); ok && ownerRef != nil {
				if name, ok := ownerRef["name"].(string); ok {
					ownerRefName = name
				}
			}
		}

		replicaInfo := types.ReplicaInfo{
			Name:         metadata["name"].(string),
			NodeID:       spec["nodeID"].(string),
			Active:       spec["active"].(bool),
			EngineName:   spec["engineName"].(string),
			CurrentState: status["currentState"].(string),
			Started:      status["started"].(bool),
			OwnerRefName: ownerRefName,
		}
		replicaInfos = append(replicaInfos, replicaInfo)
	}
	return replicaInfos, nil
}
