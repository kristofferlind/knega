package main

import (
  "path"
  "io/ioutil"
)

func codeQuality(application Application, repository Repository) error {
  reportPath := path.Join(repository.path, ".generated/code-quality")
  createIfNotExists(reportPath)

  scanCommand := "docker run --rm -v " + application.path + ":/code "
  scanCommand += "-v /var/run/docker.sock:/var/run/docker.sock "
  scanCommand += "-v /tmp/cc:/tmp/cc "
  scanCommand += "-v " + reportPath + ":/report "
  scanCommand += "-e CODECLIMATE_CODE=" + application.path + " "
  scanCommand += "codeclimate/codeclimate analyze -f html"

  result := executeCommand(scanCommand, application.repository.path)

  fileData := []byte(result)
  writeError := ioutil.WriteFile(reportPath + "/" + application.name + ".html", fileData, 0777)

  if writeError != nil {
    return writeError
  }

  return nil
}
