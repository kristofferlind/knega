package main

import (
  "log"
  "fmt"
  "strings"
  "os"
  "os/exec"

  "github.com/urfave/cli"
)

func printCommand(command *exec.Cmd) {
  fmt.Printf("==> Executing: %s\n", strings.Join(command.Args, " "))
}

func printError(err error) {
  if err != nil {
    os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    fmt.Printf("==> Output: %s\n", string(outs))
  }
}

func build(c *cli.Context) error {
  command := exec.Command("docker", "build", "-t", "app", ".")
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr

  err := command.Run()

  if err != nil {

    return err
  }

  return nil
}

func main() {
  app := cli.NewApp()
  app.Name = "Knega"
  app.Usage = "A collection of tasks for analyzing, testing, building and deploying your application"

  app.Commands = []cli.Command{
    {
      Name:  "build",
      Usage: "Builds application (uses dockerfile if it exists, otherwise tries herokuish)",
      Action: build,
    },
  }

  // start our application
  err := app.Run(os.Args)
  if err != nil {
      log.Fatal(err)
  }
}
