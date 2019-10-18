package main

import (
  "strings"
  "os/exec"
  "log"
)

func gitCloneRepository(repositoryURL string, clonePath string, directory string) string {
  return executeCommand("git clone " + repositoryURL + " " + clonePath, directory)
}

func gitPull(directory string) string {
  output := executeCommand("git pull -r", directory)

  return output
}

func gitCommit(commitMessage string, directory string) string {
  output := executeCommand("git add -A", directory)
  output += executeCommand("git commit -m \"" + commitMessage + "\"", directory)

  return output
}

func gitPush(directory string) error {
  output := gitPull(directory)
  command := exec.Command("git", "push")
  command.Dir = directory
  result, err := command.Output()

  if err != nil {
    return err
  }

  output += string(result[:])

  log.Print(output)

  return nil
}

func gitLsFiles(directory string, pattern string) []string {
  output := executeCommand("git -c core.quotepath=off ls-files " + pattern, directory)
  relativePaths := strings.Split(output, "\n")

  return relativePaths
}
