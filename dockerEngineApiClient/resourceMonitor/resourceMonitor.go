package resourceMonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ResourceMonitor(apiClient *client.Client, ctx context.Context, containerIDOrName string, interval time.Duration) {
	// Keep running it no condition for loop important for continuous monitoring
	for {
		containerStats, err := apiClient.ContainerStats(ctx, containerIDOrName, true)
		if err != nil {
			panic(err)
		}
		defer containerStats.Body.Close()

		var stats types.Stats
		if err := json.NewDecoder(containerStats.Body).Decode(&stats); err != nil {
			panic(err)
		}
		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
		cpuUsagePercent := 0.0
		if systemDelta > 0.0 {
			cpuUsagePercent = (cpuDelta / systemDelta) * float64(len(stats.PreCPUStats.CPUUsage.PercpuUsage)) * 100.0
		}
		fmt.Printf("Container %s - CPU Usage: %.2f%%, Memory Usage: %d bytes\n", containerIDOrName, cpuUsagePercent, stats.MemoryStats.Usage)

		// runs after the fixed interval
		time.Sleep(interval)
	}
}
