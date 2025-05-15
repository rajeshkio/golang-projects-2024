package engine

import (
	"context"
	"encoding/json"
	"fmt"

	types "github.com/rk280392/harvesterNavigator/internal/models"
	"k8s.io/client-go/kubernetes"
)

// FindEngineDetails retrieves engine information associated with a specific volume
func FindEngineDetails(client *kubernetes.Clientset, volumeName, absPath, namespace, resource string) ([]types.EngineInfo, error) {
	// Get all engines
	engines, err := client.RESTClient().Get().
		AbsPath(absPath).
		Namespace(namespace).
		Resource(resource).
		Do(context.Background()).Raw()

	if err != nil {
		return nil, fmt.Errorf("failed to get engines: %s", err)
	}

	var engineData map[string]interface{}
	err = json.Unmarshal(engines, &engineData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal engines data: %s", err)
	}

	items, ok := engineData["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get items from engines response")
	}

	var engineInfos []types.EngineInfo

	// Loop through all engines and find the ones for our volume
	for _, item := range items {
		engine, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Get spec to check if this engine is for our volume
		spec, ok := engine["spec"].(map[string]interface{})
		if !ok {
			continue
		}

		// Check if this engine is for our volume
		engineVolumeName, ok := spec["volumeName"].(string)
		if !ok || engineVolumeName != volumeName {
			continue
		}

		// Get metadata for engine name
		metadata, ok := engine["metadata"].(map[string]interface{})
		if !ok {
			continue
		}

		engineName, ok := metadata["name"].(string)
		if !ok {
			continue
		}

		// Create a new engine info
		engineInfo := types.EngineInfo{
			Name:      engineName,
			Snapshots: make(map[string]types.SnapshotInfo),
		}

		// Extract active status
		if active, ok := spec["active"].(bool); ok {
			engineInfo.Active = active
		}

		// Extract nodeID
		if nodeID, ok := spec["nodeID"].(string); ok {
			engineInfo.NodeID = nodeID
		}

		// Get status info
		status, ok := engine["status"].(map[string]interface{})
		if ok {
			// Extract currentState
			if state, ok := status["currentState"].(string); ok {
				engineInfo.CurrentState = state
			}

			// Extract started status
			if started, ok := status["started"].(bool); ok {
				engineInfo.Started = started
			}

			// Extract snapshots map
			if snapshots, ok := status["snapshots"].(map[string]interface{}); ok {
				for snapID, snapData := range snapshots {
					snapshot, ok := snapData.(map[string]interface{})
					if !ok {
						continue
					}

					// Initialize snapshot info
					snapshotInfo := types.SnapshotInfo{
						Name:     snapID,
						Children: make(map[string]bool),
						Labels:   make(map[string]string),
					}

					// Extract parent
					if parent, ok := snapshot["parent"].(string); ok {
						snapshotInfo.Parent = parent
					}

					// Extract created timestamp
					if created, ok := snapshot["created"].(string); ok {
						snapshotInfo.Created = created
					}

					// Extract size
					if size, ok := snapshot["size"].(string); ok {
						snapshotInfo.Size = size
					}

					// Extract userCreated flag
					if userCreated, ok := snapshot["usercreated"].(bool); ok {
						snapshotInfo.UserCreated = userCreated
					}

					// Extract removed flag
					if removed, ok := snapshot["removed"].(bool); ok {
						snapshotInfo.Removed = removed
					}

					// Extract children map
					if children, ok := snapshot["children"].(map[string]interface{}); ok {
						for child, val := range children {
							if boolVal, ok := val.(bool); ok && boolVal {
								snapshotInfo.Children[child] = true
							}
						}
					}

					// Extract labels map
					if labels, ok := snapshot["labels"].(map[string]interface{}); ok {
						for key, val := range labels {
							if strVal, ok := val.(string); ok {
								snapshotInfo.Labels[key] = strVal
							}
						}
					}

					// Add snapshot to engine's snapshots map
					engineInfo.Snapshots[snapID] = snapshotInfo
				}
			}
		}

		// Add engine to the results
		engineInfos = append(engineInfos, engineInfo)
	}

	if len(engineInfos) == 0 {
		return nil, fmt.Errorf("no engines found for volume: %s", volumeName)
	}

	return engineInfos, nil
}
