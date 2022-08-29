# sda - Sandstorm Dependency Analysis

More often than not software systems are displayed as layers.
The top-most components use the ones below but not vice-versa.
Ideally the package/namespace/module structure reflects the layers, eg:
* features built on top of other features
* REST APIs built on top of DB repositories

However it is easy to loose track of the layering in large projects, violate the structure
and in the worst case introduce cyclic dependencies.

The _sda_ tools visualized the dependencies within a code base in a graph
and validates the cycle-freeness as part of the test pipeline.

![dependencies of top level packages in this project](https://raw.githubusercontent.com/sandstorm/dependency-analysis/main/images/sda-dependencies.svg)

## Quick Start

See the project on [github.com/sandstorm/dependency-analysis](https://github.com/sandstorm/dependency-analysis) for more information.

```shell
# search and print cycles (exit code equals the number of cycles)
sda validate src/
```

## Supported languages

* Golang
* Groovy
* Java
* PHP (when using classes)

# User documentation

The given source directory is searched for known file types.
All source units (classes in Java/PHP, files in golang) are parsed to collect all dependencies.
Source units and dependencies are described using the full name including the package.
The longest package prefix shared between all source units is then used as root package.

## Example: root package recognition

Given the following Java classes:
* `de.sandstorm.examples.sda.Main`
* `de.sandstorm.examples.sda.http.WebServer`
* `de.sandstorm.examples.sda.model.WebResponse`

The root package is `de.sandstorm.examples.sda` and the top-level packages are:
* `Main`
* `http`
* `model`

Now the dependencies between this top-level packages are analyzed.
You can look deeper into the package structure by
* adjusting the source path
* by setting `--depth` to a value higher `1`

## Validate Command

You can use the _validate_ command in the CI pipeline as a code-quality test.
The exit code equals the numer of cylces in the dependency graph (hopefully zero).

Stdout contains all cycles.

### Output in project with cycles

```
$ sda validate src/main/java
found 2 cycles:

 ┌▶ http
 |   ▼
 |  services
 └───┘

 ┌▶ commands
 |   ▼
 |  http
 |   ▼
 |  services
 └───┘
```

### Output in project without cycles

```
$ sda validate
no cycles found, everything is all right
```

## Visualize Command

This is currently out of scope of the Docker images.
See [github.com/sandstorm/dependency-analysis](https://github.com/sandstorm/dependency-analysis) for more information.

## Limitations

### Code needs to be neatly formatted

For the sake on simplicity the code parsers built on regular expressions.
They only work on normal formatted code.

### Imports only

The parsers only look at imports/usages within the files.
References without imports are ignored.

### Commented lines

The parsers ignore comments.
Commented imports are treated as normal imports.

### Unused imports

Unused imports are treated as normal imports.
