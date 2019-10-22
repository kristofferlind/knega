package main

import (
  "log"
  "fmt"
  "os"

  "github.com/spf13/viper"
)

type Repository struct {
  path string
  searchDirectories []string
  searchDepth int
  applications []Application
  baseChartPath string
  helm struct{
    username string
    password string
    repository string
    repositoryGitURL string
  }
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

  // TODO: do update as part of chart upload instead, make it a job that runs for one of the applications once all workers are done
  repository.helm.repositoryGitURL = viper.GetString("Output.HelmChart.repositoryGitURL")
  repository.helm.repository = viper.GetString("Output.HelmChart.repository")
  repository.helm.username = os.Getenv("KNEGA_HELM_USERNAME")
  repository.helm.password = os.Getenv("KNEGA_HELM_PASSWORD")

  if shouldIncludeApplications {
    directoriesExist(repository.searchDirectories)
    repository.applications = getApplications(repository)
  }

  // log.Printf("Initialized repository %s", repository.path)

  return repository
}

func getApplications(repository Repository) []Application {
  var results []Application
  log.Print("Generating hashes of applications..")
  for _, searchFolder := range repository.searchDirectories {
    applicationConfigPaths := findSubDirectoriesWithFile(searchFolder, ".app.toml", repository.searchDepth)
    for _, applicationConfigPath := range applicationConfigPaths {
      application := initializeApplication(applicationConfigPath)
      results = append(results, application)
      // log.Printf("%s: %s", application.name, application.inputsHash)
    }
  }
  return results
}
