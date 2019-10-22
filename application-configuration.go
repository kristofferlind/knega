package main

import (
  "log"
  "bytes"
  "os"

  "github.com/spf13/viper"
)

type ApplicationConfiguration struct {
  name string
  inputs struct {
    filePaths []string
    gitFilePaths []string
  }
  environment struct {
    name string
    url string
  }
  outputs struct {
    dockerImage struct {
      idFile string
      repository string
      tag string
      usernameEnv string
      passwordEnv string
    }
    helmChart struct {
      chartPath string
      packageFilePath string
      packageFileName string
      repository string
      usernameEnv string
      passwordEnv string
      repositoryGitURL string
    }
  }
  commands struct {
    check []string
    build []string
    analyze []string
    release []string
  }
}

var defaultApplicationConfiguration = []byte(`
[Input]
[Input.GitFiles]
  paths = ["*"]
`)


func getApplicationConfiguration(configurationPath string, repository Repository) ApplicationConfiguration {
  configurationFile := viper.New()
  configurationFile.SetConfigType("toml")
  configurationFile.ReadConfig(bytes.NewBuffer(defaultApplicationConfiguration))

  hasRootConfig := fileExists(repository.path + "/.knega.root.toml")

  if hasRootConfig {
    configurationFile.SetConfigName(".knega.root")
    configurationFile.AddConfigPath(repository.path)
    err := configurationFile.MergeInConfig()
    if err != nil {
      log.Fatal(err)
    }
  }

  hasApplicationConfig := fileExists(configurationPath + "/.app.toml")

  if hasApplicationConfig {
    configurationFile.SetConfigName(".app")
    configurationFile.AddConfigPath(configurationPath)
    err := configurationFile.MergeInConfig()
    if err != nil {
      log.Fatal(err)
    }
  }

  configuration := ApplicationConfiguration{
    name: configurationFile.GetString("name"),
  }

  configuration.inputs.filePaths = configurationFile.GetStringSlice("Input.Files.paths")
  configuration.inputs.gitFilePaths = configurationFile.GetStringSlice("Input.GitFiles.paths")

  configuration.outputs.dockerImage.idFile = configurationFile.GetString("Output.DockerImage.idFile")
  configuration.outputs.dockerImage.repository = configurationFile.GetString("Output.DockerImage.repository")
  configuration.outputs.dockerImage.tag = configurationFile.GetString("Output.DockerImage.tag")
  configuration.outputs.dockerImage.usernameEnv = "KNEGA_DOCKER_USERNAME" //configurationFile.GetString("Output.DockerImage.usernameEnv")
  configuration.outputs.dockerImage.passwordEnv = "KNEGA_DOCKER_PASSWORD" //configurationFile.GetString("Output.DockerImage.passwordEnv")


  configuration.outputs.helmChart.chartPath = configurationFile.GetString("Output.HelmChart.chartPath")
  configuration.outputs.helmChart.packageFilePath = configurationFile.GetString("Output.HelmChart.packageFilePath")
  configuration.outputs.helmChart.packageFileName = configurationFile.GetString("Output.HelmChart.packageFileName")
  configuration.outputs.helmChart.repository = configurationFile.GetString("Output.HelmChart.repository")
  configuration.outputs.helmChart.usernameEnv = "KNEGA_HELM_USERNAME" //configurationFile.GetString("Output.HelmChart.usernameEnv")
  configuration.outputs.helmChart.passwordEnv = "KNEGA_HELM_PASSWORD" //configurationFile.GetString("Output.HelmChart.passwordEnv")
  configuration.outputs.helmChart.repositoryGitURL = configurationFile.GetString("Output.HelmChart.repositoryGitURL")

  configuration.commands.check = configurationFile.GetStringSlice("Check.commands")
  configuration.commands.build = configurationFile.GetStringSlice("Build.commands")
  configuration.commands.analyze = configurationFile.GetStringSlice("Analyze.commands")
  configuration.commands.release = configurationFile.GetStringSlice("Release.commands")

  configuration.environment.name = os.Getenv("KNEGA_ENVIRONMENT")
  if configuration.environment.name != "" {
    configuration.environment.url = configurationFile.GetString(configuration.environment.name + ".url")
  }

  return configuration
}
