package main

import (
  "path"
  "os"
  "log"
  "io"
  "time"
  "math/rand"

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
    repositoryURL := "git@github.com:Modity/HelmRepository.git" //repository.outputs.chart.path
    log.Print(gitCloneRepository(repositoryURL, repositoryName, generatedPath))
  }
  packageFileName := application.name + "-1.0.0-" + application.inputsHash + ".tgz"
  packagePath := path.Join(".generated/", packageFileName)
  packageDestinationPath := path.Join(helmRepositoryPath, packageFileName)
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


func init() {
	rand.Seed(time.Now().UnixNano())
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}
