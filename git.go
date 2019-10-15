package main

func gitLsFiles(directory string, pattern string) []string {
  result, err := exec.command("git", "-c", "core.quotepath=off", "ls-files", pattern).Directory(directory).Run()
  if err != nil {
    log.Fatal(err)
  }
  if result.exitCode != 0 {
    log.Fatal("Something went wrong when running git ls-files")
  }
  output := result.StrOutput()
  relativePaths := strings.Split(output, "\n")
}
