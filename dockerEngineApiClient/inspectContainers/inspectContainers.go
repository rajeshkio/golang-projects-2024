package inspectContainers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func InspectContainers(containerIDOrName string) (types.ContainerJSON, error) {

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	apiClient.NegotiateAPIVersion(ctx)
	defer apiClient.Close()

	containerInfo, err := apiClient.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		return types.ContainerJSON{}, err
	}

	return containerInfo, nil

}
