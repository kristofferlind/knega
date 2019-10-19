package main

import (
  "log"
  "fmt"

  "github.com/spf13/viper"
  "github.com/urfave/cli"
)

type Repository struct {
  path string
  searchDirectories []string
  searchDepth int
  // commitId string
  // isWorkTreeDirty bool
  applications []Application
  baseChartPath string
  helmRepositoryCloneURL string
  dockerRepository string
  dockerRepositoryUsername string
  dockerRepositoryPassword string
}

func initializeRepository(cliContext *cli.Context, shouldIncludeApplications bool) Repository {
  workingDirectory := getWorkingDirectory()
  repositoryPath := findParentDirectoryWithFile(workingDirectory, ".knega.root.toml")

  viper.SetConfigName(".knega.root")
  viper.AddConfigPath(repositoryPath)
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
  }

  repository := Repository{
    path: repositoryPath,
    baseChartPath: repositoryPath + viper.GetString("baseChartPath"),
    searchDirectories: viper.GetStringSlice("applicationPaths"),
    searchDepth: viper.GetInt("applicationSearchDepth"),
    helmRepositoryCloneURL: viper.GetString("helmRepositoryCloneURL"),
    dockerRepository: viper.GetString("dockerRepository"),
    dockerRepositoryUsername: cliContext.String("docker-username"),
    dockerRepositoryPassword: cliContext.String("docker-password"),
  }

  if shouldIncludeApplications {
    directoriesExist(repository.searchDirectories)
    repository.applications = getApplications(cliContext, repository)
  }

  // log.Printf("Initialized repository %s", repository.path)

  return repository
}

func getApplications(cliContext *cli.Context, repository Repository) []Application {
  var results []Application
  for _, searchFolder := range repository.searchDirectories {
    applicationConfigPaths := findSubDirectoriesWithFile(searchFolder, ".app.toml", repository.searchDepth)
    for _, applicationConfigPath := range applicationConfigPaths {
      application := initializeApplication(cliContext, applicationConfigPath)
      results = append(results, application)
      log.Printf("%s: %s", application.name, application.inputsHash)
    }
  }
  return results
}
