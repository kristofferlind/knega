package main

import (
  "strings"
  "log"
  "path/filepath"
)

// type Action func(cli.Context) error

type Job struct {
  application Application
  action string
  commands []string
  workingDirectory string
}

// TOOD: use injectVariables instead
func injectCommandVariables(rawValue string, application Application, repository Repository) string {
  value := strings.Replace(rawValue, "$ROOT", repository.path, -1)
  value = strings.Replace(value, "$INPUTS_HASH", application.inputsHash, -1)
  value = strings.Replace(value, "$APPLICATION_NAME", application.name, -1)

  value = filepath.Clean(value)

  return value
}

// TODO: use injectVariablesArray instead
func injectCommandVariablesArray(rawValues []string, application Application, repository Repository) []string {
  var values []string
  for _, rawValue := range rawValues {
    value := injectCommandVariables(rawValue, application, repository)
    values = append(values, value)
  }

  return values
}

func createJob(repository Repository, application Application, action string) Job {
  var commands []string
  switch action {
    case "check":
      commands = application.commands.check
    case "build":
      commands = application.commands.build
    case "analyze":
      commands = application.commands.analyze
    case "release":
      commands = application.commands.release
    default:
      log.Fatal("Action is not available")
  }

  commands = injectCommandVariablesArray(commands, application, repository)
  job := Job{
    action: action,
    application: application,
    commands: commands,
    workingDirectory: application.path,
  }

  return job
}
