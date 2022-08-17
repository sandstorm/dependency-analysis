# WIP: Sandstorm Dependency Analysis

* move to github in future
* analyses and visualizes dependencies between
    * Java

# User documentation

See `sda help`

# Developer documentation

* [install Golang](https://golangdocs.com/install-go-mac-os)
* `brew instal graphviz`
## Start from source

```sh
go run .
# or
go run . helloWorld
# or
go test ./parsing
```

## How to add a new command

```sh
# install COBRA CLI (see links to docs below)
cobra-cli add helloWorld
# and adjust generated files
```

## Helpful development commands

```sh
# auto-format all .go file
find . -type f -name '*.go' | xargs -Ifile go fmt file
# or
go fmt ./analysis

# print docs
go doc -all ./dataStructures
# or
go doc ./dataStructures

# docs in browser
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:8080
```

## Documentation of Libraries and Tools

* [Command Line Interface Library: COBRA](https://github.com/spf13/cobra)
* [Printing the Graph: GraphViz Examples](https://renenyffenegger.ch/notes/tools/Graphviz/examples/index)
* [GraphViz rank=same: placing node on the same level](https://stackoverflow.com/questions/14879617/layering-in-graphviz-using-dot-language)
* [colorgrad: Go (Golang) color scales library for data visualization, charts, games, maps, generative art and others](https://github.com/mazznoer/colorgrad)
* [VS Code: Markdown All in One](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one)

# Glossary

* source-unit - smallest source module, eg classes in PHP and Java, prototypes in Fusion
* package - location of a source-unit in a hierarchy, eg packages ind Java, namespaces in PHP, folders in JavaScript
* package segment - one step in the package hierarchy, e.g. _de.sandstorm.test_ consists of the three segments _[de sandstorm test]_
* root package - largest package prefix shared between all source-units
* node - box in a graph
* edge - arrow in a graph
