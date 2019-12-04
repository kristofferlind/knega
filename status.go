package main

import (
  "fmt"
  "os"
  "text/tabwriter"
)

var (
  Info = Teal
  Success = Green
  Warning = Yellow
  Fatal = Red
)

var (
  Red     = Color("\033[1;31m%s\033[0m")
  Yellow  = Color("\033[1;33m%s\033[0m")
  Green   = Color("\033[1;32m%s\033[0m")
  Teal    = Color("\033[1;36m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
  sprint := func(args ...interface{}) string {
    return fmt.Sprintf(colorString,
      fmt.Sprint(args...))
  }
  return sprint
}

func printBuildStatus(applications []Application) {
  writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)
  for _, application := range applications {
    if application.changeStatus == Unknown {
      fmt.Fprintf(writer, "%s: \t " + Warning("Unknown") + "\n", application.name)
    }
    if application.changeStatus == Dirty {
      fmt.Fprintf(writer, "%s: \t " + Fatal("Rebuild required") + "\n", application.name)
    }
    if application.changeStatus == Pristine {
      fmt.Fprintf(writer, "%s: \t " + Success("Artifacts found") + "\n", application.name)
    }
  }
  writer.Flush()
}
