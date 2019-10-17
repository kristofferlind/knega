package main

import (
  "log"
  "fmt"

  "github.com/spf13/viper"
)

type Repository struct {
  path string
  searchDirectories []string
  searchDepth int
  // commitId string
  // isWorkTreeDirty bool
  applications []Application
  baseChartPath string
}

func initializeRepository(shouldIncludeApplications bool) Repository {
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
  }

  if shouldIncludeApplications {
    directoriesExist(repository.searchDirectories)
    repository.applications = getApplications(repository)
  }

  log.Printf("Initialized repository %s", repository.path)

  return repository
}

func getApplications(repository Repository) []Application {
  var results []Application
  for _, searchFolder := range repository.searchDirectories {
    applicationConfigPaths := findSubDirectoriesWithFile(searchFolder, ".app.toml", repository.searchDepth)
    for _, applicationConfigPath := range applicationConfigPaths {
      application := initializeApplication(applicationConfigPath)
      results = append(results, application)
      log.Printf("%s: %s", application.name, application.inputsHash)
    }
  }
  return results
}
