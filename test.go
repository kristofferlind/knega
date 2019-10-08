package main

import (
  "os"
  "os/exec"
  "fmt"

  "github.com/urfave/cli"
)

func test(c *cli.Context) error {
  workingDirectory, wdErr := os.Getwd()
  if wdErr != nil {
    fmt.Println(wdErr)
  }
  command := exec.Command(
    "docker", "run",
    "-e", "BUILDPACK_URL",
    "-v", workingDirectory + ":/tmp/app:ro",
    "gliderlabs/herokuish", "/bin/herokuish", "buildpack", "test",
  )
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  err := command.Run()

  if err != nil {

    return err
  }

  return nil
}
