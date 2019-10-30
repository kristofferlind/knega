package main

import (
  "log"
  "context"
  "encoding/json"
  "encoding/base64"
  "os"

  "github.com/urfave/cli"
	"docker.io/go-docker"
  "docker.io/go-docker/api/types"
)

func dockerUpload(cliContext *cli.Context, application Application) error {
  idFile := "container.id"
  containerId := readFile(idFile)
  username := application.docker.username
  password := application.docker.password
  dockerRepository := application.docker.repository
  tag := application.docker.tag
  fullTag := dockerRepository + ":" + tag

  log.Printf("----------------------------->ContainerID: %s", containerId)

  executeCommand("docker login -u " + username + " -p " + password + " " + dockerRepository, application.path)
  executeCommand("docker image tag " + containerId + " " + fullTag, application.path)
  executeCommand("docker push " + fullTag, application.path)

  return nil
}

func dockerImageExists(imageName string, imageTag string, application *Application) bool {
  context := context.Background()
  cli, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

  inspectAuthConfig := types.AuthConfig{
		Username: application.docker.username,
    Password: application.docker.password,
	}
	encodedJSON, err := json.Marshal(inspectAuthConfig)
	if err != nil {
		log.Fatal(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

  fullImagePath := application.docker.repository + ":" + imageTag
  _, inspectErr := cli.DistributionInspect(context, fullImagePath, authStr)
  if inspectErr != nil {
    log.Print(inspectErr)
    return false
  }
  // log.Printf("%s: Found existing docker image", application.name)

  return true
}

func dockerVulnerabilityScan(cliContext *cli.Context, application Application) error {
  // imageName := application.name + ":" + application.inputsHash
  idFile := "container.id"
  containerId := readFile(idFile)

  // generate report
  generatedPath := application.repository.path + "/.generated"
  analyzePath := generatedPath + "/analyze"
  if !directoryExists(generatedPath) {
    os.Mkdir(generatedPath, 0777)
  }
  if !directoryExists(analyzePath) {
    os.Mkdir(analyzePath, 0777)
  }
  reportPath := analyzePath + "/" + application.name + ".json"
  executeCommand("trivy --no-progress --exit-code 1 -f json -o " + reportPath + " " + containerId, application.path)

  return nil
}
