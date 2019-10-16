package main

import (
  "github.com/spf13/viper"
)

type ApplicationConfiguration struct {
  name string
  fileInputPatterns []string
  gitFileInputPatterns []string
  outputs []BuildOutput
  commands struct {
    check string
    build string
    analyze string
    release string
  }
}

type BuildOutput struct {
  outputType string
}

/*
  should replace:
  - $APPLICATION_NAME
  - $ROOT
  - $INPUT_HASH (needs to be done in phase after its been calculated)
*/
func injectVariables(value string) string {
  //TODO: replace with all variables that can be used in configuration files
  return value
}

func injectVariablesArray(values []string) []string {
  var result []string
  for _, value := range values {
    transformedValue := injectVariables(value)
    result = append(result, transformedValue)
  }

  return result
}

// TODO: can override functions be just one that takes multiple types?
func handleStringSliceOverride(rootValue []string, applicationValue []string) []string {
  if applicationValue != nil {
    return applicationValue
  } else {
    return rootValue
  }
}

func handleStringOverride(rootValue string, applicationValue string) string {
  if applicationValue != "" {
    return applicationValue
  } else {
    return rootValue
  }
}

func getApplicationConfiguration(configurationPath string) ApplicationConfiguration {
  configurationFile := viper.New()
  configurationFile.SetConfigName(".app")
  configurationFile.AddConfigPath(configurationPath)

  rootFileInputPatterns := configurationFile.GetStringSlice("Input.Files.patterns")
  applicationFileInputPatterns := configurationFile.GetStringSlice("Input.Files.patterns")
  rawFileInputPatterns := handleStringSliceOverride(rootFileInputPatterns, applicationFileInputPatterns)
  fileInputPatterns := injectVariablesArray(rawFileInputPatterns)

  rootGitInputPaths := viper.GetStringSlice("Input.GitFiles.paths")
  applicationGitInputPaths := configurationFile.GetStringSlice("Input.GitFiles.paths")
  rawGitFileInputPatterns := handleStringSliceOverride(rootGitInputPaths, applicationGitInputPaths)
  gitFileInputPatterns := injectVariablesArray(rawGitFileInputPatterns)

  rootCheckCommand := viper.GetString("Check.command")
  applicationCheckCommand := configurationFile.GetString("Check.command")
  checkCommand := handleStringOverride(rootCheckCommand, applicationCheckCommand)

  rootBuildCommand := viper.GetString("Build.command")
  applicationBuildCommand := configurationFile.GetString("Build.command")
  buildCommand := handleStringOverride(rootBuildCommand, applicationBuildCommand)

  rootAnalyzeCommand := viper.GetString("Analyze.command")
  applicationAnalyzeCommand := configurationFile.GetString("Analyze.command")
  analyzeCommand := handleStringOverride(rootAnalyzeCommand, applicationAnalyzeCommand)

  rootReleaseCommand := viper.GetString("Release.command")
  applicationReleaseCommand := configurationFile.GetString("Release.command")
  releaseCommand := handleStringOverride(rootReleaseCommand, applicationReleaseCommand)

  configuration := ApplicationConfiguration{
    name: configurationFile.GetString("name"),
    fileInputPatterns: fileInputPatterns,
    gitFileInputPatterns: gitFileInputPatterns,
    // outputs: []BuildOutput,
  }

  configuration.commands.check = checkCommand // or default knega test
  configuration.commands.build = buildCommand // or default knega build
  configuration.commands.analyze = analyzeCommand // or default knega analyze (codequality, performance, dependency vulnerabilities, docker image vulnerabilities)
  configuration.commands.release = releaseCommand // or default knega release

  return configuration
}
