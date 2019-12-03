package main

import (
  "sync"
  "log"
  "time"
)

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
