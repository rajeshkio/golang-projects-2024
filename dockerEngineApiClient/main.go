package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <container-id-or-name>")
		os.Exit(1)
	}
	containerIDOrName := os.Args[1]

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err)
	}
	apiClient.NegotiateAPIVersion(ctx)
	defer apiClient.Close()

	containers, err := apiClient.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container ID: %s\n", containers.ID)
	fmt.Printf("Container Image: %s\n", containers.Image)
	fmt.Printf("Container Name: %s\n", containers.Name)
	fmt.Printf("Container Status: %s\n", containers.State.Status)
	fmt.Printf("Container Ports: %s\n", containers.NetworkSettings.Ports)
	fmt.Printf("Container Env: %s\n", containers.Config.Env)
}
