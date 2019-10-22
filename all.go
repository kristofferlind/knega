package main

import (
  "runtime"
  "log"
  "time"
  "sync"
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

func all(action string) error {
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

func getChangedApplications(applications []Application) []Application {
  var changedApplications []Application
  var asyncChangedChecks sync.WaitGroup

  changedApplicationsChannel := make(chan Application, len(applications))

  log.Print("Checking if docker images/helm charts exist with those hashes")

  for _, application := range applications {
    asyncChangedChecks.Add(1)
    go func(changedApplicationsChannel chan<- Application, application Application) {
      var hasChanges = application.hasChanges()
      if hasChanges {
        changedApplicationsChannel<-application
      }
      asyncChangedChecks.Done()
    }(changedApplicationsChannel, application)
  }
  asyncChangedChecks.Wait()
  close(changedApplicationsChannel)
  for changedApplication := range changedApplicationsChannel {
    changedApplications = append(changedApplications, changedApplication)
  }

  return changedApplications
}

func printStatus(applications []Application) {
  for _, application := range applications {
    log.Printf("%s requires rebuild", application.name)
  }
  if len(applications) == 0 {
    log.Print("No rebuilds required")
  }
}

func changed(action string) error {
  startTime := time.Now()
  repository := initializeRepository(true)
  var jobs []Job

  changedApplications := getChangedApplications(repository.applications)
  printStatus(changedApplications)

  if len(changedApplications) > 0 {
    for _, application := range changedApplications {
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
  }

  endTime := time.Now()
  timeTaken := endTime.Sub(startTime)
  log.Printf("Total time taken: %s", timeTaken)

  return nil
}
