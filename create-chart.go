package main

import (
  "os"
  "os/exec"
  "log"

  "github.com/urfave/cli"
  "github.com/otiai10/copy"
  "github.com/spf13/viper"
)

func createChart(c *cli.Context, repository Repository) error {
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

  // load
  defaultValues := viper.New()
  defaultValues.SetConfigName("values")
  defaultValues.AddConfigPath("./.generated/default-app") // change to new chart name

  // TODO: not yet available, enable once merged (https://github.com/spf13/viper/pull/635)
  // what to do as a workaround until then?
  // lets try rewriting all files to lowercase for now, done at the end
  // defaultValues.SetKeysCaseSensitive(true)

  err := defaultValues.ReadInConfig()
  if err != nil {
    log.Fatal(err)
  }
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

  version := c.String("application-version")
  if version != "" {
    chart.Set("version", version)
  }

  if appConfig.IsSet(("ingressUrl")) {
    ingressUrl := appConfig.GetString("ingress.url")
    defaultValues.Set("ingress.enabled", true)  // set from deploy-values
    defaultValues.Set("ingress.url", ingressUrl)
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

  /* All values set in deploy-values.yaml needs to overwrite those values in defaultValues, then skip that file in artifacts */
  if fileExists("deploy-values.yml") {
    err := copy.Copy("deploy-values.yml", ".generated/deploy-values.yml")
    if err != nil {
      log.Fatal(err)
    }
  } else {
    os.Create(".generated/deploy-values.yml")
  }

  return nil
}



/*
upload chart

replicate scripts/upload.sh
*/
