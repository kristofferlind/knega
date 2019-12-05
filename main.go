package main

import (
  "log"
  "os"
  "path"

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

  /*
    TODO: if currentDirectory has appconfig, initializeApplication
    if rootConfig is found initialize repository at application.repository
  */

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
      Name:  "release",
      Usage: "Deploy",
      Action: func(cliContext *cli.Context) error {
        repository := initializeRepository(false)
        var application Application
        currentWorkingDirectory := getWorkingDirectory()

        if (path.Clean(currentWorkingDirectory) != path.Clean(repository.path)) {
          application = initializeApplication(currentWorkingDirectory)
        }

        return release(cliContext, application)
      },
    },
    {
      Name: "chart",
      Usage: "chart <action> handles chart actions like create and upload",
      Subcommands: []cli.Command{
        {
          Name: "create",
          Usage: "knega chart create, creates chart based on app configuration",
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)
            var application Application
            currentWorkingDirectory := getWorkingDirectory()

            if (path.Clean(currentWorkingDirectory) != path.Clean(repository.path)) {
              application = initializeApplication(currentWorkingDirectory)
            }

            return createChart(cliContext, application, repository)
          },
        },
        {
          Name: "upload",
          Usage: "knega chart upload, uploads chart to helm repository (only works for git based repositories currently)",
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)
            var application Application
            currentWorkingDirectory := getWorkingDirectory()

            if (path.Clean(currentWorkingDirectory) != path.Clean(repository.path)) {
              application = initializeApplication(currentWorkingDirectory)
            }

            return uploadChart(cliContext, application)
          },
        },
        {
          Name: "setup-repository",
          Usage: "knega chart setup-repository, needed for helm exist check to work (done seperately because it's slow)",
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)

            return setupHelmRepository(cliContext, repository)
          },
        },
        {
          Name: "update-index",
          Usage: "knega chart update-index, updates repository index (done seperately to avoid conflicts while pushing",
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)

            return updateHelmIndex(cliContext, repository)
          },
        },
      },
    },
    {
      Name: "docker",
      Usage: "docker <action> handles docker actions like create and upload",
      Subcommands: []cli.Command{
        {
          Name:  "create",
          Usage: "Builds application (uses dockerfile if it exists, otherwise tries herokuish)",
          Action: build,
        },
        {
          Name: "upload",
          Usage: "knega docker upload, tags and uploads docker image to repository",
          Flags: []cli.Flag{
            cli.StringFlag{
              Name: "docker-username",
              Value: "",
              Usage: "Username used by docker login",
              EnvVar: "KNEGA_DOCKER_USERNAME",
            },
            cli.StringFlag{
              Name: "docker-password",
              Value: "",
              Usage: "Password used by docker login",
              EnvVar: "KNEGA_DOCKER_PASSWORD",
            },
          },
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)
            var application Application
            currentWorkingDirectory := getWorkingDirectory()

            if (path.Clean(currentWorkingDirectory) != path.Clean(repository.path)) {
              application = initializeApplication(currentWorkingDirectory)
            }

            return dockerUpload(cliContext, application)
          },
        },
        {
          Name: "vulnerability-scan",
          Usage: "knega docker vulnerability-scan, scans docker image for known vulnerabilities (trivy)",
          Flags: []cli.Flag{
            cli.StringFlag{
              Name: "exit-code",
              Value: "1",
              Usage: "Set exit-code for vulnerabilities found",
            },
          },
          Action: func(cliContext *cli.Context) error {
            repository := initializeRepository(false)
            var application Application
            currentWorkingDirectory := getWorkingDirectory()

            if (path.Clean(currentWorkingDirectory) != path.Clean(repository.path)) {
              application = initializeApplication(currentWorkingDirectory)
            }

            return dockerVulnerabilityScan(cliContext, application)
          },
        },
      },
    },
    {
      Name: "all",
      Usage: "all <action> will run action for all applications",
      Subcommands: []cli.Command{
        {
          Name:  "check",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(c *cli.Context) error {
            return all("check")
          },
        },
        {
          Name:  "build",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return all("build")
          },
        },
        {
          Name:  "analyze",
          Usage: "Run analyze command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return all("analyze")
          },
        },
        {
          Name:  "release",
          Usage: "Run release command in all applications, passing in $INPUTS_HASH",
          Action: func(context *cli.Context) error {
            return all("release")
          },
        },
      },
    },
    {
      Name: "changed",
      Usage: "changed <action> will run action for all applications with changes",
      Flags: []cli.Flag{
        cli.StringFlag{
          Name: "docker-username",
          Value: "",
          Usage: "Username used by docker login",
          EnvVar: "KNEGA_DOCKER_USERNAME",
        },
        cli.StringFlag{
          Name: "docker-password",
          Value: "",
          Usage: "Password used by docker login",
          EnvVar: "KNEGA_DOCKER_PASSWORD",
        },
      },
      Subcommands: []cli.Command{
        {
          Name:  "check",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(c *cli.Context) error {
            return changed("check")
          },
        },
        {
          Name:  "build",
          Usage: "Run build command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return changed("build")
          },
        },
        {
          Name:  "analyze",
          Usage: "Run analyze command defined in application configs where changes have occurred",
          Action: func(context *cli.Context) error {
            return changed("analyze")
          },
        },
        {
          Name:  "release",
          Usage: "Run release command in all applications, passing in $INPUTS_HASH",
          Action: func(context *cli.Context) error {
            return changed("release")
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
