package main

import (
  "log"
	"context"

  "github.com/urfave/cli"
	"docker.io/go-docker"
  "docker.io/go-docker/api/types"
  "github.com/prune998/docker-registry-client/registry"
)

func dockerUpload(cliContext *cli.Context, application Application, repository Repository) error {
  idFile := "container.id"
  containerId := readFile(idFile)
  username := cliContext.String("docker-username")
  password := cliContext.String("docker-password")
  dockerRepository := repository.dockerRepository
  appRepository := dockerRepository + "/" + application.name
  tag := application.inputsHash
  fullTag := appRepository + ":" + tag

  log.Printf("----------------------------->ContainerID: %s", containerId)

  executeCommand("docker login -u " + username + " -p " + password + " " + dockerRepository, application.path)
  executeCommand("docker image tag " + containerId + " " + fullTag, application.path)
  executeCommand("docker push " + fullTag, application.path)

  return nil
}

func dockerImageExists(imageName string, imageTag string, application *Application) bool {
  username := ""
  password := ""
  registryURL := application.repository.dockerRepository + "/" + imageName

  hub, err := registry.New(registryURL, username, password, log.Printf)
	if err != nil {
		log.Fatal("error connecting to hub, %v", err)
	}

	tags, err := hub.Tags(imageName)
	if err != nil {
		return false
	}
	for _, value := range tags {
		if value == imageTag {
      return true
    }
	}

  return false
}

func dockerImageExistsGoDocker(imageName string, imageTag string, application *Application) bool {
  context := context.Background()
  cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	authConfig := types.AuthConfig{
		Username: "",
    Password: "",
    Auth: "basic",
    ServerAddress: application.repository.dockerRepository, // + "/" + application.name,
  }

  loginResult, loginErr := cli.RegistryLogin(context, authConfig)
  if loginErr != nil {
    log.Fatal(loginErr)
  }
  log.Print(loginResult)

  searchOptions := types.ImageSearchOptions{
    Limit: 10,
  }

  searchResults, searchError := cli.ImageSearch(context, "", searchOptions)
  if searchError != nil {
    log.Fatal(searchError)
  }
  log.Print(searchResults)

  inspect, inspectErr := cli.DistributionInspect(context, "image", "")
  if inspectErr != nil {
    log.Fatal(inspectErr)
  }
  log.Print(inspect)

  return false
}
