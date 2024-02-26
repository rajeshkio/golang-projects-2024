package checkLogs

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CheckLogs(apiClient *client.Client, ctx context.Context, containerIDOrName string, follow bool, timeStamp string, logLevel string) error {

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Follow:     follow,
		Timestamps: true,
		Since:      timeStamp,
	}

	if logLevel != "" {
		if logLevel == "error" {
			options.ShowStdout = true
			options.ShowStderr = true
		} else if logLevel == "info" {
			options.ShowStdout = true
			options.ShowStderr = false
		} else {
			fmt.Println("Only info and error are correct loglevels")
		}
	}

	// Follow container logs
	containerLogs, err := apiClient.ContainerLogs(ctx, containerIDOrName, options)
	if err != nil {
		return err
	}
	defer containerLogs.Close()

	if follow {
		return followLogs(ctx, containerLogs)
	}
	// If follow is false, print logs and return
	_, err = io.Copy(os.Stdout, containerLogs)
	if err != nil {
		return err
	}

	return nil

}

func followLogs(ctx context.Context, containerLogs io.ReadCloser) error {
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
