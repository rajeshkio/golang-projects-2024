package listAndReadFiles

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ReadContentsOfFile(apiClient *client.Client, ctx context.Context, containerIDOrName string, filePath string) (string, error) {

	// Read the contents of filepath
	readContent, _, err := apiClient.CopyFromContainer(ctx, containerIDOrName, filePath)
	if err != nil {
		return "", err
	}
	defer readContent.Close()
	content, err := io.ReadAll(readContent)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ListFilesAndDirectories(apiClient *client.Client, ctx context.Context, containerIDOrName string, path string) ([]string, error) {

	// ecec into the container and run ls -l to get the list of files and directories
	execId, err := apiClient.ContainerExecCreate(ctx, containerIDOrName, types.ExecConfig{
		Cmd:          []string{"ls", "-l", path},
		AttachStderr: true,
		AttachStdout: true,
	})
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.ContainerExecAttach(ctx, execId.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	scanner := bufio.NewScanner(resp.Reader)
	var files []string
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // splits a string into substrings based on white space (spaces, tabs, or newlines)
		if len(fields) > 8 {
			files = append(files, fields[8]) // field 8 contains the file name eg drwxr-xr-x 1 root root    4096 Sep 20 16:43 nginx
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

