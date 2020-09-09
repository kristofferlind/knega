package main

import (
  "sync"
  "log"
  "time"
)

func checkApplications(applications []Application) []Application {
  var checkedApplications []Application
  var asyncChangedChecks sync.WaitGroup

  checkedApplicationsChannel := make(chan Application, len(applications))

  log.Print("Checking if artifacts exist with those hashes")

  for _, application := range applications {
    asyncChangedChecks.Add(1)
    go func(checkedApplicationsChannel chan<- Application, application Application) {
      application.hasChanges()
      checkedApplicationsChannel<-application
      asyncChangedChecks.Done()
    }(checkedApplicationsChannel, application)
  }
  asyncChangedChecks.Wait()
  close(checkedApplicationsChannel)
  for checkedApplication := range checkedApplicationsChannel {
    checkedApplications = append(checkedApplications, checkedApplication)
  }

  return checkedApplications
}

func changed(action string) error {
  startTime := time.Now()
  repository := initializeRepository(true)
  var jobs []Job

  checkedApplications := checkApplications(repository.applications)
  printBuildStatus(checkedApplications)
  var changedApplications []Application
  for _, application := range checkedApplications {
    if application.changeStatus != Pristine {
      changedApplications = append(changedApplications, application)
    }
  }

  if len(changedApplications) > 0 {
    for _, application := range changedApplications {
      job := createJob(repository, application, action)
      jobs = append(jobs, job)
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
