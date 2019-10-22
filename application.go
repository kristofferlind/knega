package main

import (
  "path"
  "path/filepath"
  "log"
  "crypto/sha512"
  "encoding/hex"
  "strings"
  "os"
)

type Application struct {
  path string
  name string
  inputsHash string
  environment struct {
    name string
    url string
  }
  docker struct {
    idFile string
    repository string
    tag string
    username string
    password string
  }
  helm struct {
    chartPath string
    packageFilePath string
    packageFileName string
    repository string
    username string
    password string
    repositoryGitURL string
  }
  commands struct {
    check []string
    build []string
    analyze []string
    release []string
  }
  repository Repository
}

func initializeApplication(applicationPath string) Application {
  shouldInitializeApplications := false
  repository := initializeRepository(shouldInitializeApplications)

  applicationConfiguration := getApplicationConfiguration(applicationPath, repository)

  application := Application{
    name: applicationConfiguration.name,
    path: applicationPath,
    repository: repository,
  }

  inputPaths := resolveInputPaths(applicationConfiguration.inputs.filePaths, applicationConfiguration.inputs.gitFilePaths, application)

  applicationInputs := initializeBuildInputs(inputPaths)
  application.inputsHash = generateInputsHash(applicationInputs)

  application.docker.idFile = injectVariables(applicationConfiguration.outputs.dockerImage.idFile, application)
  application.docker.repository = injectVariables(applicationConfiguration.outputs.dockerImage.repository, application)
  application.docker.tag = injectVariables(applicationConfiguration.outputs.dockerImage.tag, application)
  dockerUsernameEnv := injectVariables(applicationConfiguration.outputs.dockerImage.usernameEnv, application)
  dockerPasswordEnv := injectVariables(applicationConfiguration.outputs.dockerImage.passwordEnv, application)
  application.docker.username = os.Getenv(dockerUsernameEnv)
  application.docker.password = os.Getenv(dockerPasswordEnv)

  application.helm.chartPath = injectVariables(applicationConfiguration.outputs.helmChart.chartPath, application)
  application.helm.packageFilePath = injectVariables(applicationConfiguration.outputs.helmChart.packageFilePath, application)
  application.helm.packageFileName = injectVariables(applicationConfiguration.outputs.helmChart.packageFileName, application)
  application.helm.repository = injectVariables(applicationConfiguration.outputs.helmChart.repository, application)
  application.helm.repositoryGitURL = injectVariables(applicationConfiguration.outputs.helmChart.repositoryGitURL, application)
  helmUsernameEnv := injectVariables(applicationConfiguration.outputs.helmChart.usernameEnv, application)
  helmPasswordEnv := injectVariables(applicationConfiguration.outputs.helmChart.passwordEnv, application)
  application.helm.username = os.Getenv(helmUsernameEnv)
  application.helm.password = os.Getenv(helmPasswordEnv)

  application.commands.check = injectVariablesArray(applicationConfiguration.commands.check, application)
  application.commands.build = injectVariablesArray(applicationConfiguration.commands.build, application)
  application.commands.analyze = injectVariablesArray(applicationConfiguration.commands.analyze, application)
  application.commands.release = injectVariablesArray(applicationConfiguration.commands.release, application)

  application.environment.name = applicationConfiguration.environment.name
  application.environment.url = applicationConfiguration.environment.url

  return application
}

func injectVariables(rawValue string, application Application) string {
  value := strings.Replace(rawValue, "$ROOT", application.repository.path, -1)
  value = strings.Replace(value, "$INPUTS_HASH", application.inputsHash, -1)
  value = strings.Replace(value, "$APPLICATION_NAME", application.name, -1)

  if fileExists(value) {
    value = filepath.Clean(value)
  }

  return value
}

func injectVariablesArray(rawValues []string, application Application) []string {
  var values []string
  for _, rawValue := range rawValues {
    value := injectVariables(rawValue, application)
    values = append(values, value)
  }

  return values
}

func resolveInputPaths(rawFilePaths []string, rawGitFilePaths []string, application Application) []string {
  filePaths := injectVariablesArray(rawFilePaths, application)
  gitFilePaths := injectVariablesArray(rawGitFilePaths, application)

  var rawInputPaths []string
  fileInputPaths := fileInputPatternsToPaths(filePaths, application.path)
  rawInputPaths = append(rawInputPaths, fileInputPaths...)

  gitFileInputPaths := gitFileInputPatternsToPaths(gitFilePaths, application.path)
  rawInputPaths = append(rawInputPaths, gitFileInputPaths...)

  inputPaths := deDuplicateStringSlice(rawInputPaths)

  return inputPaths
}

func fileInputPatternsToPaths(patterns []string, applicationPath string) []string {
  var results []string

  for _, pattern := range patterns {
    if !isAbsolutePath(pattern) {
      pattern = path.Join(applicationPath)
    }
    matches, err := filepath.Glob(pattern)
    if err != nil {
      log.Fatal(err)
    }
    results = append(results, matches...)
  }
  return results
}

func gitFileInputPatternsToPaths(patterns []string, applicationPath string) []string {
  var results []string

  for _, pattern := range patterns {
    matches := gitLsFiles(applicationPath, pattern)
    for _, relativePath := range matches {
      path := filepath.Join(applicationPath, relativePath)
      if fileExists(path) {
        results = append(results, path)
      }
    }
  }

  return results
}

func deDuplicateStringSlice(paths []string) []string {
  pathMap := make(map[string]string)
  var results []string
  for _, path := range paths {
    exists := pathMap[path]
    if exists == "" {
      pathMap[path] = path
      results = append(results, path)
    }
  }

  return results
}

func generateInputsHash (inputs []BuildInput) string {
  hash := sha512.New()
  for _, input := range inputs {
    hash.Write(input.Hash())
  }
  checksum := hash.Sum(nil)

  return hex.EncodeToString(checksum)
}

func (application *Application) hasChanges() bool {
  packageName := application.name
  packageVersion := "1.0.0-" + application.inputsHash
  hasHelmPackage := helmPackageExists(packageName, packageVersion, application)

  imageName := application.name
  imageTag := application.inputsHash
  hasDockerImage := dockerImageExists(imageName, imageTag, application)

  return !(hasHelmPackage && hasDockerImage)
}
