package main

import (
  "os"
  "os/exec"
  "fmt"

  "github.com/urfave/cli"
)

func build(c *cli.Context) error {
  if fileExists("Dockerfile") {
    command := exec.Command("docker", "build", "-t", "app", ".")
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    err := command.Run()

    if err != nil {

      return err
    }
  } else {
    workingDirectory, wdErr := os.Getwd()
    if wdErr != nil {
      fmt.Println(wdErr)
    }
    command := exec.Command(
      "docker", "run",
      "-e", "BUILDPACK_URL",
      "-v", workingDirectory + ":/tmp/app:ro",
      "gliderlabs/herokuish", "/bin/herokuish", "buildpack", "build",
    )
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    err := command.Run()

    if err != nil {

      return err
    }
  }

  return nil
}
