# Knega
CLI tool for commonly used generic CI/CD steps (analyze, build and deploy actions)

## Inspiration
### Individual tasks
I've had a collection of bash scripts for various generic steps for a while, I bundled them in a docker image to simplify reuse a bit but it felt a bit clumsy and had issues with needing docker dind to work (which comes with a fair bit of security issues on build servers). I wanted to learn Go and thought it would be nice to have all those tasks available in a single binary, while also enabling making them work on platforms other than linux. I haven't tested it on any other platform yet, but I'm moving in that direction atleast.

### All <action> subcommands
Heavily inspired by [Baur](https://github.com/simplesurance/baur), loved the simplicity of it compared to other monorepo build tools. Most notable differences are running application pipelines in parallell and checks for changes relying on repository checks (docker image repository/helm repository) rather than a database. It also includes Check, Analyze and Release commands, with Release being run for all applications and passing in $INPUTS_HASH, which was the simplest solution I could think of for getting rollbacks to work.

Database being removed means you lose out on the build statistics (grafana dashboard).

If what Baur does works for you, use that instead. This is my first Go project, I've currently got about 10 hours of Go experience, the tool is hardly tested and have pretty crappy error handling.

## License
Since it's kind of a derivative work of [Baur](https://github.com/simplesurance/baur) I've also put the same license on it.
