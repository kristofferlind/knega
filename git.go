package main

import (
  "log"
  "os/exec"
  "strings"
)

func gitLsFiles(directory string, pattern string) []string {
  command := exec.Command("git", "-c", "core.quotepath=off", "ls-files", pattern)
  command.Dir = directory
  result, err := command.Output()
  if err != nil {
    log.Fatal(err)
  }
  output := string(result[:])
  relativePaths := strings.Split(output, "\n")

  return relativePaths
}
