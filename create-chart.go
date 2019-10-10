package main

import (
  "os"
  "os/exec"
  "log"

  "github.com/urfave/cli"
  "github.com/otiai10/copy"
)

func createChart(c *cli.Context) error {
  if ! directoryExists(".generated") {
    os.Mkdir(".generated", 0777)
  }

  chartPath := ".generated/default-app"
  // allow defining custom chart path in .knega.app.toml
  if directoryExists("chart") {
    copy.Copy("chart", chartPath)
  } else {
    // set directory in root .knega.root.toml
    baseChartPath := "../../../shared/default-app"
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

  // TODO: write known config to values.yml

  command := exec.Command("helm", "package", chartPath, "--destination", ".generated")
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  commandError := command.Run()
  if commandError != nil {
    log.Fatal(commandError)
  }

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
