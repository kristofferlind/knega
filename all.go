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
    log.Printf("Worker %d: Start work on %s", workerId, job.application.name)
    startTime := time.Now()

    result := ""

    for _, jobCommand := range job.commands {
      output := executeCommand(jobCommand, job.application.path)
      result += output
    }

    endTime := time.Now()
    timeTaken := endTime.Sub(startTime)
    log.Printf("Worker %d: Successfully ran %s for %s in %s", workerId, job.action, job.application.name, timeTaken)

    results <- result
  }
  log.Printf("Worker %d: idle, no more jobs to pick", workerId)
  asyncWorkers.Done()
}

func createWorkerPipeline (jobs []Job, resultsChannel chan<- string) {
  workerCount := runtime.NumCPU()
  jobsCount := len(jobs)
  jobsChannel := make(chan Job, jobsCount)

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
  log.Printf("Started %d workers", workerCount)
  asyncWorkers.Wait()
  close(resultsChannel)
}

func pipelineResults(results <-chan string, done chan<- bool, startTime time.Time, totalJobs int, isVerbose bool) {
  completedJobs := 0
  for result := range results {
    if isVerbose {
      log.Print(result)
    }
    completedJobs++
    currentTime := time.Now()
    log.Printf("Completed %d of %d jobs in %s", completedJobs, totalJobs, currentTime.Sub(startTime))
  }

  done <- true
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

  results := make(chan string, len(jobs))

  go pipelineResults(results, done, startTime, len(jobs), false)
  createWorkerPipeline(jobs, results)

  <-done

  endTime := time.Now()
  timeTaken := endTime.Sub(startTime)
  log.Printf("Total time taken: %s", timeTaken)

  return nil
}

func changed(c *cli.Context, action string) error {
  startTime := time.Now()
  repository := initializeRepository(true)
  var jobs []Job
  for _, application := range repository.applications {
    if application.hasChanges() {
      job := createJob(repository, application, action)
      jobs = append(jobs, job)
    }
  }

  done := make(chan bool, 1)
  defer close(done)

  results := make(chan string, len(jobs))

  go pipelineResults(results, done, startTime, len(jobs), false)
  createWorkerPipeline(jobs, results)

  <-done

  endTime := time.Now()
  timeTaken := endTime.Sub(startTime)
  log.Printf("Total time taken: %s", timeTaken)

  return nil
}
