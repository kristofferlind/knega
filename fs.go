package main

import (
  "os"
  "log"
  "fmt"
  "path"
  "path/filepath"
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
      return absolutePath
    }
  }
  log.Fatal("Could not find parent directory with file")

  // never gets here, stops lint error complaint
  return ""
}
