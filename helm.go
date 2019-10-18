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
    repositoryURL := "git@github.com:Modity/HelmRepository.git" //repository.outputs.chart.path
    log.Print(gitCloneRepository(repositoryURL, repositoryName, generatedPath))
  }

  commitMessage := "re-index"
  log.Print(createHelmRepositoryIndex(helmRepositoryPath))
  log.Print(gitCommit(commitMessage, helmRepositoryPath))
  err := gitPush(helmRepositoryPath)

  if err != nil {
    return err
  }

  return nil
}
