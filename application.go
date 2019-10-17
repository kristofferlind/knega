package main

import (
  "path"
  "path/filepath"
  "log"
  "crypto/sha512"
  "encoding/hex"
)

type Application struct {
  path string
  name string
  inputsHash string
  inputs []BuildInput
  outputs []BuildOutput
  commands struct {
    check []string
    build []string
    analyze []string
    release []string
  }
}

func initializeApplication(applicationPath string) Application {
  applicationConfiguration := getApplicationConfiguration(applicationPath)

  application := Application{
    name: applicationConfiguration.name,
    path: applicationPath,
  }

  var rawInputPaths []string
  fileInputPaths := fileInputPatternsToPaths(applicationConfiguration.fileInputPatterns, applicationPath)
  rawInputPaths = append(rawInputPaths, fileInputPaths...)

  gitFileInputPaths := gitFileInputPatternsToPaths(applicationConfiguration.gitFileInputPatterns, applicationPath)
  rawInputPaths = append(rawInputPaths, gitFileInputPaths...)

  inputPaths := deDuplicateStringSlice(rawInputPaths)
  application.inputs = initializeBuildInputs(inputPaths)
  application.inputsHash = generateInputsHash(application.inputs)

  // application.outputArtifacts :=

  return application
}

// TODO: support ** in patterns (https://github.com/simplesurance/baur/blob/master/resolve/glob/glob.go)
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
      // TODO: applicationPath is an absolute path, should be rootRelative?
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
  // check if output artifacts exist (where uploaded) for that specific hash
  return true
}
