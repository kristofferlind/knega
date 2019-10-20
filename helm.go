package main

import (
  "path"
  "os"
  "log"

  "github.com/urfave/cli"
)

func createHelmRepositoryIndex(directory string) string {
  return executeCommand("helm repo index .", directory)
}

func updateHelmIndex(cliContext *cli.Context, repository Repository) error {
  generatedPath := path.Join(repository.path, ".generated")
  if ! directoryExists(generatedPath) {
    os.Mkdir(generatedPath, 0777)
  }
  repositoryName := "helm-repo"
  helmRepositoryPath := path.Join(repository.path, ".generated/", repositoryName)
  if ! directoryExists(helmRepositoryPath) {
    repositoryURL := repository.helm.repositoryGitURL
    log.Print(gitCloneRepository(repositoryURL, repositoryName, generatedPath))
  }

  commitMessage := "re-index"
  log.Print(createHelmRepositoryIndex(helmRepositoryPath))
  log.Print(gitCommit(commitMessage, helmRepositoryPath))

  // TODO: need a retry here.. probably for the clone aswell
  err := gitPush(helmRepositoryPath)

  if err != nil {
    return err
  }

  return nil
}

func helmPackageExists(packageName string, packageVersion string, application *Application) bool {
  addRepoCommand := "helm repo add --username " + application.helm.username
  addRepoCommand += " --password " + application.helm.password
  addRepoCommand += " knega-repo " + application.helm.repository
  executeCommand(addRepoCommand, application.path)

  executeCommand("helm repo update", application.path)

  searchCommand := "helm search --version 1.0.0-" + application.inputsHash
  searchCommand += " knega-repo/" + application.name
  result := executeCommand(searchCommand, application.path)

  if result == "No results found" {
    return false
  }

  log.Printf("Found existing chart in these results: %s", result)

  return true
}
