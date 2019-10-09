package main

import (
  "os"
  "os/exec"
  // "fmt"

  "github.com/urfave/cli"
)

/*
  Just a copy of build action for now, should do the following:
  - check if Taskfile.yml exists, if it does, run build
  - otherwise run defaults from root?
*/
func pipeline(c *cli.Context) error {
  if fileExists("Taskfile.yml") {
    command := exec.Command("task", "build")
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    err := command.Run()

    if err != nil {

      return err
    }
  } else {
    // copy base Taskfile.yml?
    // task build
  }

  return nil
}
