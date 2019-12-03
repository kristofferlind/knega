package main

import (
  "github.com/urfave/cli"
)

type Filterer struct {
  tags []string
  applicationName string
}

func initializeFilterer(cliContext *cli.Context) Filterer {
  filterer := Filterer{
    applicationName: cliContext.String("application-name"),
    tags: cliContext.StringSlice("tags"),
  }

  return filterer
}

func filterApplications(applications []Application, test func(Application) bool) []Application {
  var result []Application
  for _, application := range applications {
    if test(application) {
      result = append(result, application)
    }
  }
  return result
}

func arrayAll(collection []string, test func(string) bool) bool {
  for _, collectionValue := range collection {
    if !test(collectionValue) {
      return false
    }
  }
  return true
}

func findIndex(collection []string, value string) int {
  for index, collectionValue := range collection {
    if collectionValue == value {
      return index
    }
  }
  return -1
}

func hasString(collection []string, value string) bool {
  return findIndex(collection, value) >= 0
}

func (filterer *Filterer) filter(applications []Application) []Application {
  hasTags := len(filterer.tags) > 0
  var result []Application

  if hasTags {
    hasTagsTest := func(application Application) bool {
      return arrayAll(filterer.tags, application.hasTag)
    }

    return filterApplications(applications, hasTagsTest)
  }

  return result
}
