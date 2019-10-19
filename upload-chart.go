package main

import (
  "path"
  "os"
  "log"
  "io"
  "time"

  "github.com/urfave/cli"
)

// TODO: need to pass in application
func uploadChart(context *cli.Context, application Application, repository Repository) error {
  generatedPath := path.Join(repository.path, ".generated")
  if ! directoryExists(generatedPath) {
    os.Mkdir(generatedPath, 0777)
  }

  // do one repository per application to minimze conflicts
  repositoryName := application.name + "-repo"
  helmRepositoryPath := path.Join(repository.path, ".generated/", repositoryName)
  if ! directoryExists(helmRepositoryPath) {
    repositoryURL := repository.helmRepositoryCloneURL
    log.Print(gitCloneRepository(repositoryURL, repositoryName, generatedPath))
  }
  packageFileName := application.name + "-1.0.0-" + application.inputsHash + ".tgz"
  packagePath := path.Join(".generated/", packageFileName)

  packageDestinationDirectory := path.Join(helmRepositoryPath, "charts")
  if ! directoryExists(packageDestinationDirectory) {
    os.Mkdir(packageDestinationDirectory, 0777)
  }

  packageDestinationPath := path.Join(packageDestinationDirectory, packageFileName)
  log.Print(packageDestinationPath)

  sourceFile, sourceErr := os.Open(packagePath)
  if sourceErr != nil {
    return sourceErr
  }
  defer sourceFile.Close()

  destinationFile, destinationErr := os.Create(packageDestinationPath)
  if destinationErr != nil {
    return destinationErr
  }
  defer destinationFile.Close()

  _, copyErr := io.Copy(destinationFile, sourceFile)
  if copyErr != nil {
    return copyErr
  }

  // do helm repository reindex after build instead to avoid conflicts
  // log.Print(createHelmRepositoryIndex(helmRepositoryPath))
  commitMessage := application.name + "-" + application.inputsHash
  log.Print(gitCommit(commitMessage, helmRepositoryPath))

  // TODO: retry?
  retry(10, time.Second, func() error {
    err := gitPush(helmRepositoryPath)
    if err != nil {
      return err
    }
    return nil
  })

  return nil
}
