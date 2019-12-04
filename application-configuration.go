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
    urls []string
  }
  artifacts []interface{} // will replace outputs
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
  configurationFile := getConfigurationManager(configurationPath, repository)

  configuration := ApplicationConfiguration{
    name: configurationFile.GetString("name"),
  }

  configuration.inputs.filePaths = configurationFile.GetStringSlice("Input.Files.paths")
  configuration.inputs.gitFilePaths = configurationFile.GetStringSlice("Input.GitFiles.paths")

  configuration.outputs.dockerImage.idFile = configurationFile.GetString("Output.DockerImage.idFile")
  configuration.outputs.dockerImage.repository = configurationFile.GetString("Output.DockerImage.repository")
  configuration.outputs.dockerImage.tag = configurationFile.GetString("Output.DockerImage.tag")
  configuration.outputs.dockerImage.usernameEnv = "KNEGA_DOCKER_USERNAME"
  configuration.outputs.dockerImage.passwordEnv = "KNEGA_DOCKER_PASSWORD"

  configuration.outputs.helmChart.chartPath = configurationFile.GetString("Output.HelmChart.chartPath")
  configuration.outputs.helmChart.packageFilePath = configurationFile.GetString("Output.HelmChart.packageFilePath")
  configuration.outputs.helmChart.packageFileName = configurationFile.GetString("Output.HelmChart.packageFileName")
  configuration.outputs.helmChart.repository = configurationFile.GetString("Output.HelmChart.repository")
  configuration.outputs.helmChart.usernameEnv = "KNEGA_HELM_USERNAME"
  configuration.outputs.helmChart.passwordEnv = "KNEGA_HELM_PASSWORD"
  configuration.outputs.helmChart.repositoryGitURL = configurationFile.GetString("Output.HelmChart.repositoryGitURL")

  configuration.commands.check = configurationFile.GetStringSlice("Check.commands")
  configuration.commands.build = configurationFile.GetStringSlice("Build.commands")
  configuration.commands.analyze = configurationFile.GetStringSlice("Analyze.commands")
  configuration.commands.release = configurationFile.GetStringSlice("Release.commands")

  configuration.environment.name = os.Getenv("KNEGA_ENVIRONMENT")
  if configuration.environment.name != "" {
    configuration.environment.urls = configurationFile.GetStringSlice(configuration.environment.name + ".urls")
  }

  configuration.artifacts = getArtifactsConfiguration(configurationFile)

  return configuration
}

func getConfigurationManager(configurationPath string, repository Repository) *viper.Viper {
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

  return configurationFile
}

type Artifact interface {
  exists() bool
  upload() bool
}

type DockerImageArtifact struct {
  idFile string
  repository string
  tag string
  usernameEnv string
  passwordEnv string
}

type HelmChartArtifact struct {
  chartPath string
  packageFilePath string
  packageFileName string
  repository string
  usernameEnv string
  passwordEnv string
  repositoryGitURL string
}

func getArtifactsConfiguration(configurationManager *viper.Viper) []interface{} {
  // if output defined in app config, skip all that were defined in root config
  // outputs should just be an array of artifacts
  // artifact should be an interface and each artifact type should define how to check for them
  artifacts := make([]interface{}, 2)

  dockerImage := DockerImageArtifact{
    idFile: configurationManager.GetString("Output.DockerImage.idFile"),
    repository: configurationManager.GetString("Output.DockerImage.repository"),
    tag: configurationManager.GetString("Output.DockerImage.tag"),
    usernameEnv: "KNEGA_DOCKER_USERNAME",
    passwordEnv: "KNEGA_DOCKER_PASSWORD",
  }
  artifacts[0] = dockerImage

  helmChart := HelmChartArtifact{
    chartPath: configurationManager.GetString("Output.HelmChart.chartPath"),
    packageFilePath: configurationManager.GetString("Output.HelmChart.packageFilePath"),
    packageFileName: configurationManager.GetString("Output.HelmChart.packageFileName"),
    repository: configurationManager.GetString("Output.HelmChart.repository"),
    usernameEnv: "KNEGA_HELM_USERNAME",
    passwordEnv: "KNEGA_HELM_PASSWORD",
    repositoryGitURL: configurationManager.GetString("Output.HelmChart.repositoryGitURL"),
  }
  artifacts[1] = helmChart

  return artifacts
}
