package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/rk280392/dockerEngineApiClient/checkLogs"
	"github.com/rk280392/dockerEngineApiClient/inspectContainers"
	"github.com/rk280392/dockerEngineApiClient/resourceMonitor"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: ./main <command> <container-id-or-name>")
		os.Exit(1)
	}
	ctx := context.Background()
	command := os.Args[1]
	containerIDOrName := os.Args[2]

	apiClient := initAPIClient()
	defer apiClient.Close()

	switch command {
	case "inspect":
		containerInfo, err := inspectContainers.InspectContainers(apiClient, ctx, containerIDOrName)
		if err != nil {
			fmt.Println("Failed to inspect container:", err)
			os.Exit(1)
		}
		inspectContainers.PrintContainerInfo(containerInfo)
	case "stats":
		interval := 5 * time.Second
		// By prefacing the monitorContainer function call with go, a new goroutine is spawned to execute monitorContainer concurrently with the main program.
		go resourceMonitor.ResourceMonitor(apiClient, ctx, containerIDOrName, interval)

		// Keep the main goroutine running
		select {}
	case "logs":
		err := checkLogs.CheckLogs(apiClient, ctx, containerIDOrName)
		if err != nil {
			fmt.Println("Failed to follow container logs:", err)
			os.Exit(1)
		}
	}

}

func initAPIClient() *client.Client {
	var err error
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	// Optionally negotiate API version
	apiClient.NegotiateAPIVersion(context.Background())
	return apiClient
}
