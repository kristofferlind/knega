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
}

func initializeRepository() Repository {
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
    searchDirectories: viper.GetStringSlice("applicationPaths"),
    searchDepth: viper.GetInt("applicationSearchDepth"),
  }
  directoriesExist(repository.searchDirectories)
  repository.applications = getApplications(repository)

  return repository
}

func getApplications(repository Repository) []Application {
  var results []Application
  for _, searchFolder := range repository.searchDirectories {
    applicationConfigPaths := findSubDirectoriesWithFile(searchFolder, ".app.toml", repository.searchDepth)
    for _, applicationConfigPath := range applicationConfigPaths {
      application := initializeApplication(applicationConfigPath)
      results = append(results, application)
      log.Print("Found application: " + application.path)
    }
  }
  return results
}
