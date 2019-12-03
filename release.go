package main

import (
  "os"
  "log"
  "path"
  "strings"

  "github.com/urfave/cli"
)

//TODO: move to fs
func clearDirectory(path string) {
  if directoryExists(path) {
    rmErr := os.RemoveAll(path)
    if rmErr != nil {
      log.Fatal(rmErr)
    }
  }
  os.Mkdir(path, 0777)
}

//TODO: move to fs
func createIfNotExists(path string) {
  if !directoryExists(path) {
    os.Mkdir(path, 0777)
  }
}

func release(cliContext *cli.Context, application Application) error {
  // create directories if not exist (put in root?)
  fetchPath := path.Join(application.repository.path, ".generated/charts")
  renderPath := path.Join(application.repository.path, ".generated/rendered_charts")
  createIfNotExists(fetchPath)
  createIfNotExists(renderPath)

  // helm fetch
  fetchCommand := "helm fetch --repo " + application.helm.repository
  fetchCommand += " --username " + application.helm.username
  fetchCommand += " --password " + application.helm.password
  fetchCommand += " --untar --untardir " + fetchPath
  fetchCommand += " --version 1.0.0-" + application.inputsHash
  fetchCommand += " --debug " + application.name
  executeCommand(fetchCommand, application.repository.path)

  // helm template
  renderCommand := "helm template --set environment=" + application.environment.name
  if len(application.environment.urls) > 0 {
    renderCommand += ",ingress.enabled=true,ingress.urls={" + strings.Join(application.environment.urls[:], ",") + "}"
  }

  // if len(application.environment.variables) > 0 {
  //   renderCommand += ",configMapRefs.environmentVariables=" + application.name + "-environment-variables-secret"
    /*
      set application.secretName like in auto-deploy-app (https://gitlab.com/gitlab-org/charts/auto-deploy-app)
      generate configmap like this, with an environment variable for each line:
apiVersion: v1
kind: ConfigMap
metadata:
  name: special-config
  namespace: default
data:
  SPECIAL_LEVEL: very
  SPECIAL_TYPE: charm

chart should grab variables from secret if set:
apiVersion: v1
kind: Pod
metadata:
  name: dapi-test-pod
spec:
  containers:
    - name: test-container
      image: k8s.gcr.io/busybox
      command: [ "/bin/sh", "-c", "env" ]
      envFrom:
      - configMapRef:
          name: special-config
  restartPolicy: Never
    */
    // for _, environmentVariable := range application.environment.variables {
    //
    // }
  // }
  renderCommand += " --output-dir " + renderPath
  renderCommand += " " + fetchPath + "/" + application.name
  executeCommand(renderCommand, application.repository.path)

  // should only happen once, doing it outside of knega for now
  // kubectl apply
  // releaseCommand := "kubectl apply --recursive --filename " + renderPath
  // executeCommand(releaseCommand, application.repository.path)

  return nil
}
