package main

import (
  "strings"
)

// type Action func(cli.Context) error

type Job struct {
  action string
  commands []string
  workingDirectory string
}

// TOOD: use injectVariables instead
func injectCommandVariables(rawCommand string, application Application, repository Repository) string {
  command := strings.Replace(rawCommand, "$ROOT", repository.path, -1)
  command = strings.Replace(command, "$INPUTS_HASH", application.inputsHash, -1)
  command = strings.Replace(command, "$APPLICATION_NAME", application.name, -1)

  return command
}

func injectCommandVariablesArray(rawValues []string, application Application, repository Repository) []string {
  var values []string
  for _, rawValue := range rawValues {
    value := injectCommandVariables(rawValue, application, repository)
    values = append(values, value)
  }

  return values
}

func createJob(repository Repository, application Application, action string) Job {
  rawCommands := application.commands[action]
  commands := injectCommandVariablesArray(rawCommands, application, repository)
  job := Job{
    action: action,
    commands: commands,
    workingDirectory: application.path,
  }

  return job
}
