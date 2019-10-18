package main

import (
  "runtime"
  "log"
  "time"
  "sync"

  "github.com/urfave/cli"
)

// maybe just include output when error?
func createWorker (workerId int, asyncWorkers *sync.WaitGroup, jobs <-chan Job, results chan<- string) {
  for job := range jobs {
    log.Printf("Start work on %s", job.application.name)
    startTime := time.Now()

    result := ""

    for _, jobCommand := range job.commands {
      output := executeCommand(jobCommand, job.application.path)
      result += output
    }

    endTime := time.Now()
    timeTaken := endTime.Sub(startTime)
    log.Printf("Successfully ran %s for %s in %s", job.action, job.application.name, timeTaken)

    results <- result
  }
  asyncWorkers.Done()
}

func createWorkerPipeline (done chan<- bool, jobs []Job) (<-chan string) {
  workerCount := runtime.NumCPU()
  jobsCount := len(jobs)
  jobsChannel := make(chan Job, jobsCount)
  resultsChannel := make(chan string, jobsCount)

  for _, job := range jobs {
    jobsChannel <- job
    log.Printf("Created %s job for %s", job.action, job.application.name)
  }
  close(jobsChannel)

  var asyncWorkers sync.WaitGroup
  for workerId := 1; workerId <= workerCount; workerId++ {
    asyncWorkers.Add(1)
    go createWorker(workerId, &asyncWorkers, jobsChannel, resultsChannel)
  }
  asyncWorkers.Wait()
  close(resultsChannel)
  done <- true

  log.Printf("Started %d workers", workerCount)

  return resultsChannel
}

func all(c *cli.Context, action string) error {
  startTime := time.Now()
  repository := initializeRepository(true)
  var jobs []Job
  for _, application := range repository.applications {
    job := createJob(repository, application, action)
    jobs = append(jobs, job)
  }

  done := make(chan bool, 1)
  defer close(done)

  results := createWorkerPipeline(done, jobs)

  for result := range results {
    log.Print(result)
  }

  <-done

  endTime := time.Now()
  timeTaken := endTime.Sub(startTime)
  log.Printf("Total time taken: %s", timeTaken)

  return nil
}
