package main

import (
  "github.com/urfave/cli"
)

func createChart(c *cli.Context) error {
  if directoryExists("chart") {
    // copy to .generated/chart
  } else {
    // if basechart/shared chart exists
    // generate default chart to .generated/chart
    // else
    // throw error, no app specific or shared chart defined
  }

  // generate values.yaml

  return nil
}
