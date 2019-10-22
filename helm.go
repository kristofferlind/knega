package main

import (
  "path"
  "os"
  "log"
  "strings"

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

func setupHelmRepository(cliContext *cli.Context, repository Repository) error {
  addRepoCommand := "helm repo add --username " + repository.helm.username
  addRepoCommand += " --password " + repository.helm.password
  addRepoCommand += " knega-repo " + repository.helm.repository
  executeCommand(addRepoCommand, repository.path)

  executeCommand("helm repo update", repository.path)

  return nil
}

func helmPackageExists(packageName string, packageVersion string, application *Application) bool {
  // TODO: setupHelmRepository if knega-repo does not exist or just have it run once for first application it hits
  searchCommand := "helm search --version 1.0.0-" + application.inputsHash
  searchCommand += " knega-repo/" + application.name
  result := executeCommand(searchCommand, application.path)

  if result == "No results found" {
    return false
  }

  if strings.Contains(result, "knega-repo/" + application.name) {
    log.Printf("%s: Found existing helm chart", application.name)
    return true
  }

  log.Printf("Something went wrong while checking for helm chart search command: %s returned the following results: %s", searchCommand, result)

  return false
}
