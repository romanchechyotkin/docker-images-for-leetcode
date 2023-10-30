package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func RunSubmit() bytes.Buffer {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	tar, err := archive.TarWithOptions("subtract_numbers/", &archive.TarOptions{})
	if err != nil {
		log.Fatalf("failed to tar: %v", err)
	}

	var imageName = "subtract_numbers-code"

	res, err := cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	})
	if err != nil {
		log.Fatalf("failed to build: %v", err)
	}

	var buffer bytes.Buffer
	_, err = io.Copy(io.MultiWriter(os.Stdout, &buffer), res.Body)
	if err != nil {
		fmt.Println("Failed to read build response:", err)
		os.Exit(1)
	}

	imagesID := getImageID(&buffer)
	log.Println("image built:", imagesID[len(imagesID)-1])

	createContainerCfg := &container.Config{
		Image: imageName,
	}

	var containerID string
	containerCreate, err := cli.ContainerCreate(ctx, createContainerCfg, nil, nil, nil, "test-container")
	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}
	containerID = containerCreate.ID
	log.Printf("container created: %s\n", containerID)

	err = cli.ContainerStart(ctx, containerCreate.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("failed to run container: %v", err)
	}

	logs := containerLogs(cli, ctx, containerID)

	defer remove(cli, imagesID, containerID)
	return logs
}

func getImageID(res *bytes.Buffer) []string {
	var imageID []string

	decoder := json.NewDecoder(res)
	for {
		var v map[string]interface{}
		if err := decoder.Decode(&v); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Failed to parse build response:", err)
			os.Exit(1)
		}

		if aux, ok := v["aux"]; ok {
			if auxMap, ok := aux.(map[string]interface{}); ok {
				if id, idOk := auxMap["ID"]; idOk {
					imageID = append(imageID, id.(string))
				}
			}
		}
	}
	return imageID
}

func remove(cli *client.Client, imagesID []string, containerID string) {
	fmt.Printf("-----------------------------------------\n")
	err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Println(err)
	}

	for _, image := range imagesID {
		imageRemove, err := cli.ImageRemove(context.Background(), image, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			log.Println(err)
		}
		log.Println(imageRemove)
	}
}

func containerLogs(cli *client.Client, ctx context.Context, containerID string) bytes.Buffer {
	out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		fmt.Println("failed to retrieve container logs:", err)
		os.Exit(1)
	}
	defer out.Close()

	var logsBuffer bytes.Buffer
	_, err = io.Copy(&logsBuffer, out)
	if err != nil {
		log.Println("failed to copy:", err)
	}

	return logsBuffer
}
