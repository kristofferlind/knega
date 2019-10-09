package main

import (
  "os"
)

func pathExists(path string) bool {
  info, err := os.Stat(path)
  if os.IsNotExist(err) {
      return false
  }
  return !info.IsDir()
}

func fileExists(filename string) bool {
  return pathExists(filename);
}

func directoryExists(directory string) bool {
  return pathExists(directory);
}
