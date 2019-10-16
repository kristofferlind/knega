package main

import (
  "crypto/sha512"
  "os"
  "log"
  "io"
)

type BuildInput struct {
  path string
  // absolutePath string
  // rootRelativePath string
  // applicationRelativePath string
}

func initializeBuildInputs(paths []string) []BuildInput {
  var inputs []BuildInput

  for _, path := range paths {
    input := BuildInput{
      path: path,
    }
    inputs = append(inputs, input)
  }

  return inputs
}

func (buildInput *BuildInput) Hash() []byte {
  file, err := os.Open(buildInput.path)
  if err != nil {
    log.Fatal(err)
  }

  // close file on function return
  defer file.Close()

  hash := sha512.New()

  _, copyErr := io.Copy(hash, file)
  if copyErr != nil {
    log.Fatal(copyErr)
  }

  return hash.Sum(nil)
}
