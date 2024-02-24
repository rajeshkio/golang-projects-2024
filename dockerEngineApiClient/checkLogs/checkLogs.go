package checkLogs

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CheckLogs(apiClient *client.Client, ctx context.Context, containerIDOrName string) error {

	// Follow container logs
	containerLogs, err := apiClient.ContainerLogs(ctx, containerIDOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	defer containerLogs.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			buf := make([]byte, 1024)
			n, err := containerLogs.Read(buf)
			if err != nil {
				return err
			}
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
