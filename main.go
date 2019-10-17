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

  // if currentDirectory has appconfig, initialize repository without applications
  // and initialize just the application for current directory
  // might also be a reasonable alternative to reverse the relationship between repository and applications, most actions focus on application rather than repository
  // so have applications hold a reference to repository rather than the other way around
  // then single application commands should get single application and all commands should get an array of applications
  // for now just get from repository
  repository := initializeRepository(false)

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
      Action: func(context *cli.Context) error {
        return createChart(context, repository)
      },
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
            return all(c, "check")
          },
        },
        {
          Name:  "build",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return all(context, "build")
          },
        },
        {
          Name:  "analyze",
          Usage: "Run analyze command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return all(context, "analyze")
          },
        },
        {
          Name:  "release",
          Usage: "Run release command in all applications, passing in $INPUTS_HASH",
          Action: func(context *cli.Context) error {
            return all(context, "release")
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
