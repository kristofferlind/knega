package main

import (
  "log"
  "os"

  "github.com/urfave/cli"
)

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
    {
      Name:  "test",
      Usage: "Test application (using herokuish)",
      Action: test,
    },
    {
      Name:  "pipeline",
      Usage: "Checks for Taskfile.yml, if present run that, otherwise run default",
      Action: pipeline,
    },
    {
      Name:  "create-chart",
      Usage: "Build app chart",
      Action: pipeline,
    },
  }

  err := app.Run(os.Args)
  if err != nil {
      log.Fatal(err)
  }
}
