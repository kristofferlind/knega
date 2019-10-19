package main

import (
  "os"
  "log"
  "fmt"
  "path"
  "path/filepath"
  "io/ioutil"
  "bytes"
  "strings"
)

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }

  return !info.IsDir()
}

func directoryExists(directory string) bool {
  info, err := os.Stat(directory)
  if os.IsNotExist(err) {
    return false
  }

  return info.IsDir()
}

func directoriesExist(directories []string) {
  for _, directory := range directories {
    if ! directoryExists(directory) {
      log.Fatal("Directory does not exist", directory)
    }
  }
}

func getWorkingDirectory() string {
  workingDirectory, err := os.Getwd()
  if err != nil {
    log.Fatal(fmt.Errorf("Could not get working directory path", err))
  }
  return workingDirectory
}

func findSubDirectoriesWithFile(directory string, filename string, searchDepth int) []string {
  var results []string
  globPath := ""

  for depth := 0; depth <= searchDepth; depth++ {
    searchPath := path.Join(directory, globPath, filename)
    matches, err := filepath.Glob(searchPath)
    if err != nil {
      log.Fatal(err)
    }
    for _, match := range matches {
      absolutePath, err := filepath.Abs(match)
      if err != nil {
        log.Fatal(err)
      }
      absolutePath = strings.Replace(absolutePath, filename, "", 1)
      results = append(results, absolutePath)
    }

    globPath += "*/"
  }

  return results
}

func findParentDirectoryWithFile(directory string, filename string) string {
  searchDirectoryPath := directory

  // TODO: should stop search when at system root, for now just limit depth to 10
  for depth := 0; depth <= 10; depth++ {
    searchPath := path.Join(searchDirectoryPath, filename)
    if fileExists(searchPath) {
      absolutePath, err := filepath.Abs(searchPath)
      if err != nil {
        log.Fatal(err)
      }

      directoryPath := strings.Replace(absolutePath, filename, "", 1)
      // log.Print("Found root directory: " + directoryPath)

      return directoryPath
    }
    searchDirectoryPath += "../"
  }
  log.Fatal("Could not find parent directory with file")

  // never gets here, stops lint error complaint
  return ""
}

// TODO: implement this
func isAbsolutePath(path string) bool {
  return true
}

// The horror.. need this until viper is fixed or replaced
func convertFileContentToLowerCase(path string) {
  content, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal(err)
  }
  lowerCaseContent := bytes.Replace(content, []byte(".applicationName"), []byte(".applicationname"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".replicaCount"), []byte(".replicacount"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".pullPolicy"), []byte(".pullpolicy"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".environmentVariables"), []byte(".environmentvariables"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte("ASPNETCORE_URLS"), []byte("aspnetcore_urls"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".internalPort"), []byte(".internalport"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".livenessProbe"), []byte(".livenessprobe"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".initialDelaySeconds"), []byte(".initialdelayseconds"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".timeoutSeconds"), []byte(".timeoutseconds"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".readinessProbe"), []byte(".readinessprobe"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".minReplicas"), []byte(".minreplicas"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".maxReplicas"), []byte(".maxreplicas"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".targetCPUUtilizationPercentage"), []byte("targetcpuutilizationpercentage"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".modSecurity"), []byte(".modsecurity"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".secRuleEngine"), []byte(".secruleengine"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".externalPort"), []byte(".externalport"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".podDisruptionBudget"), []byte(".poddisruptionbudget"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".minAvailable"), []byte(".minavailable"), -1)
  lowerCaseContent = bytes.Replace(lowerCaseContent, []byte(".maxUnavailable"), []byte(".maxunavailable"), -1)

  err = ioutil.WriteFile(path, lowerCaseContent, 0777)
}

func readFile(path string) string {
  output, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal(err)
  }
  return string(output[:])
}
