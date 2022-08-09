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
* [GraphViz rank=same: placing node on the same level](https://stackoverflow.com/questions/14879617/layering-in-graphviz-using-dot-language)

# Glossary

* source-unit - smallest source module, eg classes in PHP and Java, prototypes in Fusion
* package - location of a source-unit in a hierarchy, eg packages ind Java, namespaces in PHP, folders in JavaScript
* package segment - one step in the package hierarchy, e.g. _de.sandstorm.test_ consists of the three segments _[de sandstorm test]_
* root package - largest package prefix shared between all source-units
* node - box in a graph
* edge - arrow in a graph
