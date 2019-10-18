package main

import (
  "runtime"
  "log"
  "time"

  "github.com/urfave/cli"
)

// maybe just include output when error?
func createWorker (done <-chan struct{}, jobs <-chan Job, results chan<- string, errors chan<- error) {
  for job := range jobs {
    log.Printf("Start work on %s", job.application.name)
    startTime := time.Now()

    result := ""

    for _, jobCommand := range job.commands {
      output := executeCommand(jobCommand, job.application.path)
      log.Print(output)
      result += output
    }

    endTime := time.Now()
    timeTaken := endTime.Sub(startTime)
    log.Printf("Successfully ran %s for %s in %s", job.action, job.application.name, timeTaken)
    select {
      case results <- result:
      case <- done:
        return
    }
  }
}

func createWorkerPipeline (done <-chan struct{}, jobs []Job) (<-chan string, <-chan error) {
  workerCount := runtime.NumCPU()
  jobsCount := len(jobs)
  jobsChannel := make(chan Job, jobsCount)
  resultsChannel := make(chan string)
  errorsChannel := make(chan error, 1)

  for workerId := 1; workerId <= workerCount; workerId++ {
    go createWorker(done, jobsChannel, resultsChannel, errorsChannel)
  }

  log.Printf("Started %d workers", workerCount)

  for _, job := range jobs {
    jobsChannel <- job
    log.Printf("Created %s job for %s", job.action, job.application.name)
  }
  close(jobsChannel)

  return resultsChannel, errorsChannel
}

func all(c *cli.Context, action string) error {
  repository := initializeRepository(true)
  var jobs []Job
  for _, application := range repository.applications {
    job := createJob(repository, application, action)
    jobs = append(jobs, job)
  }

  done := make(chan struct{})
  defer close(done)

  results, errors := createWorkerPipeline(done, jobs)

  for result := range results {
    log.Print(result)
  }

  if err := <- errors; err != nil {
    return err
  }

  return nil
}
