# Knega
CLI tool for commonly used generic CI/CD steps (analyze, build and deploy actions)

## Usage
Tool has help sections describing all the different actions, overview found by running `knega`

### Examples
#### Monorepo actions
Example configurations can be found in [examples/](./examples), [.knega.root.toml](./examples/.knega.root.toml) should be at project root and [.app.toml](./examples/.app.toml) should exist for every application

Command | Action
--- | ---
`knega changed <action>` | runs action for applications with changes (based on hash of inputs, checks whether image & chart with that hash already exists)
`knega all <action>` | runs action for all applications
`knega changed check` | runs check commands for applications with changes
`knega changed build` | runs build commands for applications with changes
`knega changed analyze` | runs analyze commands for applications with changes
`knega all release` | runs release commands for all applications

#### Individual actions
Command | Action
--- | ---
`knega chart create` | Create chart
`knega chart upload` | Upload chart to repository
`knega docker create` | Builds based on dockerfile if exists, otherwise tries herokuish build
`knega docker upload` | Upload docker image to repository
`knega docker test` | Runs tests (only works if application matches a herokuish buildpack)

## Inspiration
### Individual actions
I've had a collection of bash scripts for various generic steps for a while, I bundled them in a docker image to simplify reuse a bit but it wasn't great and had issues with needing docker dind to work (which comes with a fair bit of security issues on build servers). I wanted to learn Go and thought it would be nice to have all those tasks available in a single binary, while also enabling making them work on platforms other than linux. I haven't tested it on any other platform yet, but I'm trying to avoid OS specific stuff.

### Monorepo actions
Heavily inspired by [Baur](https://github.com/simplesurance/baur), loved the simplicity of it compared to other monorepo build tools. Most notable differences are running application pipelines in parallell and checks for changes relying on repository checks (docker image repository/helm repository) rather than a database. It also includes Check, Analyze and Release commands and makes $ROOT, $APPLICATION_NAME and $INPUTS_HASH available so that any part can be replaced except input hash generation and existence checks (might make existence checks replaceable aswell).

Database being removed means you lose out on the build statistics (grafana dashboard).

If what Baur does works for you, use that instead. This is my first Go project, I've currently got about 20 hours of Go experience, it's not well tested (I only know that it works for my current projects) and error handling isn't great.

This does however fix the issues I had with Baur and has been a great project for learning a bit of Go.

## License
Since parts of it is kind of a derivative work of [Baur](https://github.com/simplesurance/baur) I've also put the same license on it.
