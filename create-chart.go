package main

import (
  "os"
  "os/exec"
  "log"

  "github.com/urfave/cli"
  "github.com/otiai10/copy"
  "github.com/spf13/viper"
)

func createChart(cliContext *cli.Context, application Application, repository Repository) error {
  log.Printf("application-version: %s", cliContext.String("application-version"))
  // TODO: create ensureCleanDirectoryExists func in fs/path
  if directoryExists(".generated") {
    rmErr := os.RemoveAll(".generated")
    if rmErr != nil {
      log.Fatal(rmErr)
    }
  }
  os.Mkdir(".generated", 0777)

  chartPath := ".generated/default-app"
  // allow defining custom chart path in .knega.app.toml
  if directoryExists("chart") {
    copy.Copy("chart", chartPath)
  } else {
    // set directory in root .knega.root.toml
    baseChartPath := repository.baseChartPath
    if directoryExists(baseChartPath) {
      err := copy.Copy(baseChartPath, chartPath)
      if err != nil {
        log.Fatal(err)
      }
    } else {
      // generate default chart? download gitlab auto deploy chart?
      log.Fatal("To run create-chart you need to have either an application specific chart or a base chart defined")
    }
  }

  if fileExists("deploy-values.yml") {
    err := copy.Copy("deploy-values.yml", ".generated/deploy-values.yml")
    if err != nil {
      log.Fatal(err)
    }
  } else {
    os.Create(".generated/deploy-values.yml")
  }

  // load
  defaultValues := viper.New()
  defaultValues.SetConfigName("values")
  defaultValues.AddConfigPath("./.generated/default-app") // change to new chart name
  err := defaultValues.ReadInConfig()
  if err != nil {
    log.Fatal(err)
  }

  defaultValues.SetConfigName("deploy-values")
  defaultValues.AddConfigPath(".generated")
  err = defaultValues.MergeInConfig()
  if err != nil {
    log.Fatal(err)
  }

  defaultValues.SetConfigName("values")
  defaultValues.AddConfigPath("./.generated/default-app") // change to new chart name

  // TODO: not yet available, enable once merged (https://github.com/spf13/viper/pull/635)
  // what to do as a workaround until then?
  // lets try rewriting all files to lowercase for now, done at the end
  // defaultValues.SetKeysCaseSensitive(true)

  appConfig := viper.New()
  appConfig.SetConfigName(".app")
  appConfig.AddConfigPath(".")
  err = appConfig.ReadInConfig()
  if err != nil {
    log.Fatal(err)
  }
  applicationName := appConfig.GetString("name")

  defaultValues.Set("applicationName", applicationName)

  // TODO: should set chart name and .generated/<chart-directory> to applicationName-<inputHash>
  chart := viper.New()
  chart.SetConfigName("Chart")
  chart.AddConfigPath("./.generated/default-app") // change to new chart name
  err = chart.ReadInConfig()
  if err != nil {
    log.Fatal(err)
  }
  chart.Set("name", applicationName)

  inputsHash := application.inputsHash
  chartVersion := chart.GetString("version")
  if inputsHash != "" {
    chart.Set("appVersion", inputsHash)
    chart.Set("version", chartVersion + "-" + inputsHash)
    defaultValues.Set("image.tag", inputsHash)
  }

  defaultValues.WriteConfig()
  chart.WriteConfig()

  renameErr := os.Rename("./.generated/default-app", "./.generated/" + applicationName)
  if renameErr != nil {
    log.Fatal(renameErr)
  }

  // TODO: remove these once case sensitivity can be turned on for viper (1.5?)
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/_helpers.tpl")
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/deployment.yaml")
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/hpa.yaml")
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/ingress.yaml")
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/pdb.yaml")
  convertFileContentToLowerCase(".generated/" + applicationName + "/templates/service.yaml")

  command := exec.Command("helm", "package", "./.generated/" + applicationName, "--destination", ".generated")
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  commandError := command.Run()
  if commandError != nil {
    log.Fatal(commandError)
  }

  return nil
}
