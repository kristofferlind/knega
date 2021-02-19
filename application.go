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

type ChangeStatus int

const (
  Unknown ChangeStatus = iota
  Dirty
  Pristine
)

func (changeStatus ChangeStatus) String() string {
  return [...]string{"Unknown", "Dirty", "Pristine"}[changeStatus]
}

type Application struct {
  path string
  name string
  inputsHash string
  environment struct {
    name string
    urls []string
    variables []string
  }
  changeStatus ChangeStatus
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
    repositoryOCI string
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
  application.helm.repositoryOCI = injectVariables(applicationConfiguration.outputs.helmChart.repositoryOCI, application)

  helmUsernameEnv := injectVariables(applicationConfiguration.outputs.helmChart.usernameEnv, application)
  helmPasswordEnv := injectVariables(applicationConfiguration.outputs.helmChart.passwordEnv, application)
  application.helm.username = os.Getenv(helmUsernameEnv)
  application.helm.password = os.Getenv(helmPasswordEnv)

  application.commands.check = injectVariablesArray(applicationConfiguration.commands.check, application)
  application.commands.build = injectVariablesArray(applicationConfiguration.commands.build, application)
  application.commands.analyze = injectVariablesArray(applicationConfiguration.commands.analyze, application)
  application.commands.release = injectVariablesArray(applicationConfiguration.commands.release, application)

  application.environment.name = applicationConfiguration.environment.name
  application.environment.urls = applicationConfiguration.environment.urls

  application.changeStatus = Unknown

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
  hash := sha512.New384()
  for _, input := range inputs {
    hash.Write(input.Hash())
  }
  checksum := hash.Sum(nil)

  return hex.EncodeToString(checksum)
}

func (application *Application) hasChanges() bool {
  hasHelmPackage := false
  if application.helm.repositoryOCI != "" {
    hasHelmPackage = ociHelmPackageExists(application)
  } else if application.helm.repository != "" {
    hasHelmPackage = helmPackageExists(application)
  }

  imageName := application.name
  imageTag := application.inputsHash
  hasDockerImage := dockerImageExists(imageName, imageTag, application)

  skipHelmPackage := (application.helm.repository == "" && application.helm.repositoryOCI == "")
  skipDockerImage := application.docker.repository == ""
  hasArtifacts := ((skipHelmPackage || hasHelmPackage) && (skipDockerImage || hasDockerImage))

  if hasArtifacts {
    application.changeStatus = Pristine
  } else {
    application.changeStatus = Dirty
  }

  return !hasArtifacts
}

func (application *Application) hasTag(tag string) bool {
  // check if has tag
  return false
}

func (application *Application) nameContains(value string) bool {
  // check if name contains string
  return false
}
