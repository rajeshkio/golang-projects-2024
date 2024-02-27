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
	"github.com/rk280392/dockerEngineApiClient/listAndReadFiles"
	"github.com/rk280392/dockerEngineApiClient/resourceMonitor"
)

func main() {

	command := flag.String("command", "", "Command to execute: 'inspect', 'stats', 'logs', 'list', 'read'")
	containerIDOrName := flag.String("container", "", "Container ID or name")
	follow := flag.Bool("follow", false, "Follow logs in real-time")
	timeStamp := flag.String("timestamp", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
	interval := flag.Duration("interval", 5*time.Second, "Interval for resource monitoring")
	loglevel := flag.String("loglevel", "info", "loglevel log filtering (e.g., 'info', 'error)")
	readFilePath := flag.String("readfilepath", "", "read the contents of the file")

	// Parse command-line arguments
	flag.Parse()

	if *command == "" || *containerIDOrName == "" {
		fmt.Println("Usage: ./main -command <command> -container <container-id-or-name> [-follow] [-timestamp <timestamp>] [-interval <interval>] [-loglevel <loglevel>] [-readFilePath <filepath>]")
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
		err := checkLogs.CheckLogs(apiClient, ctx, *containerIDOrName, *follow, *timeStamp, *loglevel)
		if err != nil {
			fmt.Println("Failed to follow container logs:", err)
			os.Exit(1)
		}
	case "list":
		output, err := listAndReadFiles.ListFilesAndDirectories(apiClient, ctx, *containerIDOrName, "/etc")
		if err != nil {
			fmt.Println("Failed tp hey list of files", err)
		}
		fmt.Println(output)
	case "read":
		output, err := listAndReadFiles.ReadContentsOfFile(apiClient, ctx, *containerIDOrName, *readFilePath)
		if err != nil {
			fmt.Println("Failed to read contents of the file")
		}
		fmt.Println(output)
	default:
		fmt.Println("Usage: ./program <command>")
		fmt.Println("Available commands: inspect, stats, logs, list, read")
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
