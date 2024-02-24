package inspectContainers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func InspectContainers(apiClient *client.Client, ctx context.Context, containerIDOrName string) (types.ContainerJSON, error) {

	// Get the container infor
	containerInfo, err := apiClient.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return containerInfo, nil

}

func PrintContainerInfo(containerInfo types.ContainerJSON) {
	// Assuming containerInfo is a map[string]interface{} representing container information
	fmt.Printf("Container ID: %s\n", containerInfo.ID)
	fmt.Printf("Container Image: %s\n", containerInfo.Image)
	fmt.Printf("Container Name: %s\n", containerInfo.Name)
	fmt.Printf("Container State: %s\n", containerInfo.State.Status)
	fmt.Printf("Container Ports: %v\n", containerInfo.NetworkSettings.Ports)
	fmt.Printf("Container Env: %v\n", containerInfo.Config.Env)
}
