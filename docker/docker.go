package docker

import (
	"context"
	// "fmt"
	"log"
	// "github.com/docker/docker/pkg/stack"
	"github.com/moby/moby/api/types/container"
	client "github.com/moby/moby/client"
	// Tcontainer "github.com/moby/moby/api/types/container"
)

const maxContainerId int = 12

func ListContainers(includeAllContainers bool) ([]ContainerInfo, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal("could not create api client")
		return nil, nil
	}

	options := client.ContainerListOptions{
		All: includeAllContainers,
	}

	containers, _ := apiClient.ContainerList(context.TODO(), options)
	containerInfos := make([]ContainerInfo, 0)

	for _, container := range containers.Items {
		summary := ContainerInfo{
			Id:     container.ID[0:maxContainerId],
			Name:   container.Names[0],
			State:  container.State,
			Status: container.Status,
		}

		containerInfos = append(containerInfos, summary)
	}

	return containerInfos, nil
}

func StartContainer(containerId string) ([]ContainerInfo, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal("could not create api client")
		return nil, nil
	}

	apiClient.ContainerStart(context.TODO(), containerId, client.ContainerStartOptions{})

	containers, err := ListContainers(true)

	if err != nil {
		return []ContainerInfo{}, err
	}

	return containers, nil
}

func StopContainer(containerId string) ([]ContainerInfo, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal("could not create api client")
		return nil, nil
	}

	apiClient.ContainerStop(context.TODO(), containerId, client.ContainerStopOptions{})

	containers, err := ListContainers(true)

	if err != nil {
		return []ContainerInfo{}, err
	}

	return containers, nil
}

func RestartContainer(containerId string) ([]ContainerInfo, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal("could not create api client")
		return nil, nil
	}

	apiClient.ContainerRestart(context.TODO(), containerId, client.ContainerRestartOptions{})

	containers, err := ListContainers(true)

	if err != nil {
		return []ContainerInfo{}, err
	}

	return containers, nil

}

type ContainerInfo struct {
	Id     string
	Name   string
	State  container.ContainerState
	Status string
}
