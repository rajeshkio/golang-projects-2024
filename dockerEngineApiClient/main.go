package main

import (
	"fmt"
	"os"

	"github.com/rk280392/dockerEngineApiClient/inspectContainers"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <container-id-or-name>")
		os.Exit(1)
	}
	containerIDOrName := os.Args[1]
	containerInfo, err := inspectContainers.InspectContainers(containerIDOrName)
	fmt.Printf("Container ID: %s\n", containerInfo.ID)
	fmt.Printf("Container Image: %s\n", containerInfo.Image)
	fmt.Printf("Container Name: %s\n", containerInfo.Name)
	fmt.Printf("Container Status: %s\n", containerInfo.State.Status)
	fmt.Printf("Container Ports: %s\n", containerInfo.NetworkSettings.Ports)
	fmt.Printf("Container Env: %s\n", containerInfo.Config.Env)
	if err != nil {
		fmt.Println(err)
	}
}
