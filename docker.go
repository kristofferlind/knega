package main

import (
  "log"
  "context"
  "encoding/json"
  "encoding/base64"

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

func dockerImageExistsOther(imageName string, imageTag string, application *Application) bool {
  username := ""
  password := ""
  registryURL := "https://" + application.repository.dockerRepository + "/" + imageName

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

func dockerImageExists(imageName string, imageTag string, application *Application) bool {
  context := context.Background()
  cli, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

  inspectAuthConfig := types.AuthConfig{
		Username: application.repository.dockerRepositoryUsername,
    Password: application.repository.dockerRepositoryPassword,
	}
	encodedJSON, err := json.Marshal(inspectAuthConfig)
	if err != nil {
		log.Fatal(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

  fullImagePath := application.repository.dockerRepository + "/" + imageName + ":" + imageTag
  inspect, inspectErr := cli.DistributionInspect(context, fullImagePath, authStr)
  if inspectErr != nil {
    log.Print(inspectErr)
    return false
  }
  log.Printf("Found iamge %s at %s", inspect.Descriptor.Digest, fullImagePath)

  return true
}
