# WIP: Sandstorm Dependency Analysis

* move to github in future
* analyses and visualizes dependencies between
    * Java

# User documentation

See `dca help`

# Developer documentation

* [install Golang](https://golangdocs.com/install-go-mac-os)
* `brew instal graphviz`
## Start from source

```sh
go run github.com/sandstorm/dependency-analysis
# or
go run github.com/sandstorm/dependency-analysis helloWorld
# or
go test github.com/sandstorm/dependency-analysis/parsing
```


## How to add a new command

```sh
# install COBRA CLI (see links to docs below)
cobra-cli add helloWorld
# and adjust generated files
```

## Documentation of Libraries

* [Command Line Interface Library: COBRA](https://github.com/spf13/cobra)
* [Printing the Graph: GraphViz Examples](https://renenyffenegger.ch/notes/tools/Graphviz/examples/index)
