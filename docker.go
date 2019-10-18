package main

import (
  "log"

  "github.com/urfave/cli"
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
