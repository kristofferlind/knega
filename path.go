package main

import (
  "os"
)

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }

  return !info.IsDir()
}

func directoryExists(directory string) bool {
  info, err := os.Stat(directory)
  if os.IsNotExist(err) {
    return false
  }

  return info.IsDir()
}
