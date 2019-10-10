// package main

// import (
//   "runtime"
//   "log"
//   "fmt"

//   "github.com/urfave/cli"
//   "github.com/spf13/viper"
// )

// type action func(cli.Context) error
// type job interface {
//   application application,
//   action action
// }

// type application interface {
//   config()
// }

// func worker (jobs <-chan job, results chan<- error) {
//   for job := range jobs {
//     // set working directory to application.path
//     // logs, err := execute action
//     // print logs
//     // results <- err
//   }
// }

// // func remoteWorker as above, execute as job on k8s cluster?

// func getRootConfig() {
//   viper.SetConfigName(".knega.root")
//   viper.AddConfigPath(".")
//   err := viper.ReadInConfig()
//   if err != nil {
//     log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
//   }
// }

// func initializeApplication(applicationPath string) application {
//   applicationConfig := viper.new()
//   viper.SetConfigName(".knega.app")
//   viper.AddConfigPath(applicationPath)
//   application := &application{
//     config: viperConfig
//   }

//   return application
// }

// func getAppllications() {
//   // find application configs in specified directories, searchLevels deep
// }

// func all(c *cli.Context, action *action) error {
//   // locate root config .knega.root.toml
//   getRootConfig()
//   // applications := locate all app configs below root { config, path, hasChanges, inputsHash } .knega.app.toml
//   getApplications()
//   workerCount := runtime.NUMCPU() * 2
//   remoteWorkerCount := rootConfig.remoteWorker.count
//   jobs := make(chan job, applications.length)
//   results := make(chan error, applications.length)

//   // for cpu threads, initiate worker
//   for workerId := 1; workerId <= workerCount; workerId++ {
//     go worker(jobs, results)
//   }

//   // for configured remoteWorkers count, initiate remote worker
//   // for remoteWorkerId := 1; remoteWorkerId <= remoteWorkerCount; remoteWorkerId++ {
//   //   go remoteWorker(jobs, results)
//   // }

//   // for applications create jobs and send on jobs channel, close jobs channel
//   for range applications {
//     job := new Job(application, action)
//     jobs <- job
//   }
//   close(jobs)

//   // for results <-results
//   for range applications {
//     // if error, terminate application
//     // if success save inputsHash
//     <-results
//   }

//   // successfully ran command for all changed applications
// }
