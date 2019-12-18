package main

import (
  "path"
  "strings"

  "github.com/urfave/cli"
)

func release(cliContext *cli.Context, application Application) error {
  fetchPath := path.Join(application.repository.path, ".generated/charts")
  renderPath := path.Join(application.repository.path, ".generated/rendered_charts")
  createIfNotExists(fetchPath)
  createIfNotExists(renderPath)

  commitId := getLatestCommit(application.repository)
  repositoryCommitUrl := strings.Replace(application.helm.repository, "master", commitId, 1)

  fetchCommand := "helm fetch --repo " + repositoryCommitUrl
  fetchCommand += " --username " + application.helm.username
  fetchCommand += " --password " + application.helm.password
  fetchCommand += " --untar --untardir " + fetchPath
  fetchCommand += " --version 1.0.0-" + application.inputsHash
  fetchCommand += " --debug " + application.name
  executeCommand(fetchCommand, application.repository.path)

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
  // releaseCommand := "kubectl apply --recursive --filename " + renderPath
  // executeCommand(releaseCommand, application.repository.path)

  return nil
}
