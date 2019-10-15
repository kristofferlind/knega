package main

import (
  "runtime"

  "github.com/urfave/cli"
)

type Action func(cli.Context) error

type Job struct {
  application Application
  action string
}

func worker (jobs <-chan Job, results chan<- error) {
  for job := range jobs {
    // set working directory to application.path
    // logs, err := execute action
    // print logs
    // results <- err
  }
}

// func remoteWorker as above, execute as job on k8s cluster?

func all(c *cli.Context, action string) error {
  repository := initializeRepository()
  workerCount := runtime.NumCPU() * 2
  applicationsCount := len(repository.applications)
  jobs := make(chan Job, applicationsCount)
  results := make(chan error, applicationsCount)

  // for cpu threads, initiate worker
  for workerId := 1; workerId <= workerCount; workerId++ {
    go worker(jobs, results)
  }

  // for configured remoteWorkers count, initiate remote worker
  // for remoteWorkerId := 1; remoteWorkerId <= repository.remoteWorkerCount; remoteWorkerId++ {
  //   go remoteWorker(jobs, results)
  // }

  // for applications create jobs and send on jobs channel, close jobs channel
  for _, application := range repository.applications {
    job := Job{application, action}
    jobs <- job
  }
  close(jobs)

  // for results <-results
  for range repository.applications {
    // if error, terminate application
    // if success save inputsHash, dont define outputs?, make dockerfile the default and tag those with hash?
    // that way we can skip saving in database and just check if that docker image exists locally or on connected container registry
    // maybe do the same for helm repository?
    <-results
  }

  // successfully ran command for all changed applications
  return nil
}
