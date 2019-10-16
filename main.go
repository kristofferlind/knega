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

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "application-version",
      Value: "",
      Usage: "Version tag to be used for build/create-chart commands",
    },
  }

  app.Commands = []cli.Command{
    {
      Name:  "build",
      Usage: "Builds application (uses dockerfile if it exists, otherwise tries herokuish)",
      Action: build,
    },
    {
      Name:  "test",
      Usage: "Test application (using herokuish)",
      Category: "Analyze",
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
      Action: createChart,
    },
    {
      Name:  "release",
      Usage: "Deploy",
      Action: test,
    },
    {
      Name: "all",
      Usage: "all <action> will run action for all applications with changes",
      Subcommands: []cli.Command{
        {
          Name:  "check",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(c *cli.Context) error {
            return all(c, "Build")
          },
        },
        {
          Name:  "build",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(c *cli.Context) error {
            return all(c, "Build")
          },
        },
        {
          Name:  "analyze",
          Usage: "Run analyze command defined in application configs where changes have occurred",
          Action: func(c *cli.Context) error {
            return all(c, "Build")
          },
        },
        {
          Name:  "release",
          Usage: "Run release command in all applications, passing in $INPUTS_HASH",
          Action: func(c *cli.Context) error {
            return all(c, "Build")
          },
        },
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
      log.Fatal(err)
  }
}
