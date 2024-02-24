package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/rk280392/dockerEngineApiClient/checkLogs"
	"github.com/rk280392/dockerEngineApiClient/inspectContainers"
	"github.com/rk280392/dockerEngineApiClient/resourceMonitor"
)

func main() {

	command := flag.String("command", "", "Command to execute: 'inspect', 'stats', or 'logs'")
	containerIDOrName := flag.String("container", "", "Container ID or name")
	follow := flag.Bool("follow", false, "Follow logs in real-time")
	timeStamp := flag.String("timestamp", "", "Timestamp for log filtering (e.g., '2024-02-24T16:07:41.140905466Z')")
	interval := flag.Duration("interval", 5*time.Second, "Interval for resource monitoring")

	// Parse command-line arguments
	flag.Parse()

	if *command == "" || *containerIDOrName == "" {
		fmt.Println("Usage: ./main -command <command> -container <container-id-or-name> [-follow] [-timestamp <timestamp>] [-interval <interval>]")
		os.Exit(1)
	}
	ctx := context.Background()

	apiClient := initAPIClient()
	defer apiClient.Close()

	switch *command {
	case "inspect":
		containerInfo, err := inspectContainers.InspectContainers(apiClient, ctx, *containerIDOrName)
		if err != nil {
			fmt.Println("Failed to inspect container:", err)
			os.Exit(1)
		}
		inspectContainers.PrintContainerInfo(containerInfo)
	case "stats":
		// By prefacing the monitorContainer function call with go, a new goroutine is spawned to execute monitorContainer concurrently with the main program.
		go resourceMonitor.ResourceMonitor(apiClient, ctx, *containerIDOrName, *interval)

		// Keep the main goroutine running
		select {}
	case "logs":
		err := checkLogs.CheckLogs(apiClient, ctx, *containerIDOrName, *follow, *timeStamp)
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
