package main

import (
  "log"
)

type Job struct {
  application Application
  action string
  commands []string
  workingDirectory string
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

  job := Job{
    action: action,
    application: application,
    commands: commands,
    workingDirectory: application.path,
  }

  return job
}
