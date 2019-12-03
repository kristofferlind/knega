package main

import (
  "os/exec"
  "log"
  "strings"
  "bytes"
  "bufio"
)

func executeCommand(command string, directory string) string {
  commandParts := strings.Split(command, " ")
  commandExecutor := exec.Command(commandParts[0], commandParts[1:]...)
  commandExecutor.Dir = directory

  // log.Printf("Executing command: %s in %s", command, directory)

  logReader, initErr := commandExecutor.StdoutPipe()
  if initErr != nil {
    log.Printf("Executing command: %s in %s, crashed before receiving any output", command, directory)
    log.Fatal(initErr)
  }

  commandExecutor.Stderr = commandExecutor.Stdout

  startError := commandExecutor.Start()
  
  if startError != nil {
    log.Printf("Executing command: %s in %s, crashed before receiving any output", command, directory)
    log.Fatal(startError)
  }

  var outputBuffer bytes.Buffer

  firstLine := true
  scanner := bufio.NewScanner(logReader)

  for scanner.Scan() {
    if firstLine {
      firstLine = false
    } else {
      outputBuffer.WriteRune('\n')
    }
    outputBuffer.Write(scanner.Bytes())
  }

  if err := scanner.Err(); err != nil {
    commandExecutor.Wait()
    log.Print(string(outputBuffer.Bytes()))
    log.Fatal(err)
  }

  waitError := commandExecutor.Wait()
  if waitError != nil {
    log.Print(string(outputBuffer.Bytes()))
    log.Fatal(waitError)
  }

  return string(outputBuffer.Bytes())
}
