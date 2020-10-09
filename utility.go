package main

import (
  "os"
)

func getMin(a int, b int) int {
  if a <= b {
    return a
  }
  return b
}

func IsTrace() bool {
  traceValue := os.Getenv("KNEGA_TRACE")
  if traceValue == "true" {
    return true
  }

  return false
}
