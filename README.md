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

![dependencies of top level packages in this project](./images/sda-dependencies.svg)

## Installation

You can download [the latest release](https://github.com/sandstorm/dependency-analysis/releases/) or compile the project on your own.

```shell
# download release
open https://github.com/sandstorm/dependency-analysis/releases/

# … or install via homebrew
brew install sandstorm/tap/sandstorm-dependency-analysis

# … or compile sda locally
go build -o ~/bin/sda

# install GraphViz (with Homebrew on OSX)
brew install graphviz

# show available commands
sda --help

# render and open dependency graph
sda visualize src/

# search and print cycles (exit code equals the number of cycles)
sda validate src/
```

For your CI pipeline there is an image on [hub.docker.com](https://hub.docker.com/r/sandstormmedia/dependency-analysis).

```shell
# use this image in your CI pipeline
docker run --rm sandstormmedia/dependency-analysis sda --help

# note that the image does not (yet) include GraphViz
# but you can generate the DOT graph and pipe it
# to your local GraphViz installation
docker run --rm \
  -v $(pwd):/src \
  sandstormmedia/dependency-analysis \
  sda visualize /src -o stdout \
| dot -T svg -o output.svg
```

## Supported languages

* Golang
* Groovy
* Java
* PHP (when using classes)
* JavaScript (when using modules)
* TypeScript (when using modules)

# Table of Content

- [sda - Sandstorm Dependency Analysis](#sda---sandstorm-dependency-analysis)
  - [Installation](#installation)
  - [Supported languages](#supported-languages)
- [Table of Content](#table-of-content)
- [User documentation](#user-documentation)
  - [Example: root package recognition](#example-root-package-recognition)
  - [Validate Command](#validate-command)
    - [Output in project with cycles](#output-in-project-with-cycles)
    - [Output in project without cycles](#output-in-project-without-cycles)
    - [Gitlab CI integration](#gitlab-ci-integration)
  - [Visualize Command](#visualize-command)
    - [Graph with single cycle](#graph-with-single-cycle)
    - [Graph with several cycles](#graph-with-several-cycles)
  - [How to debug and break cycles](#how-to-debug-and-break-cycles)
    - [Cycles disappear with increased depth](#cycles-disappear-with-increased-depth)
    - [Cycles between classes](#cycles-between-classes)
  - [Settings](#settings)
    - [--depth](#--depth)
    - [--include](#--include)
    - [--max-cycles (validate only)](#--max-cycles-validate-only)
    - [--graphLabel (visualize only)](#--graphlabel-visualize-only)
    - [--output (visualize only)](#--output-visualize-only)
    - [--show-image (visualize and OSX only)](#--show-image-visualize-and-osx-only)
    - [--show-node-labels (visualize only)](#--show-node-labels-visualize-only)
    - [-T, --type string (visualize only)](#-t---type-string-visualize-only)
  - [Limitations](#limitations)
    - [Code needs to be neatly formatted](#code-needs-to-be-neatly-formatted)
    - [Imports only](#imports-only)
    - [Commented lines](#commented-lines)
    - [Unused imports](#unused-imports)
- [Developer documentation](#developer-documentation)
  - [Start from source](#start-from-source)
  - [Release](#release)
  - [Add a new command](#add-a-new-command)
  - [Helpful development commands](#helpful-development-commands)
  - [Documentation of Libraries and Tools](#documentation-of-libraries-and-tools)
- [Glossary](#glossary)

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

### Gitlab CI integration

```yaml
test-code-dependencies:
  image: sandstormmedia/dependency-analysis:v1
  stage: test
  needs: [] # executes earlier
  dependencies: [] # starts faster
  script:
    - sda validate src/main/java --depth 1
    - sda validate src/main/java --depth 2
    - sda validate src/main/java --depth 3 --max-cycles 1
```

## Visualize Command

To get an overview over a legacy project or to debug and break cyclic dependencies
you can render the dependency graph into an image, such as

![dependencies of top level packages in this project](./images/sda-dependencies.svg)

The color of the nodes indicated the position in the layering:
more blue packages are higher in the layering,
more green ones lower.

The colors are particulary helpful in tangeld graphs of messy dependencies.

Note that I omit the node labels in the following examples.

### Graph with single cycle

Edges on a cycle have warm colors.

![Graph with one single cycle](./images/cyclic-dependency.svg)

### Graph with several cycles

In larger graphs with multiple cycles the color of the edges
helps to follow single cycles.
Since in the presence of cycles layers cease to exist,
the graph becomes tangled.
The colors of the nodes still indicate on which level the node
*might* be.

![Graph with several cycles](./images/several-cyclic-dependencies.svg)

## How to debug and break cycles

If you plan to clean up a project and remove cycles I suggest the following:
* start on thick edges to break multiple cycles at once
* increase the `--depth` to get some insights where the cycles come from

### Cycles disappear with increased depth

Lucky you, there are no cycles between classes.
Re-ordering a few classes and packages should do the trick.

### Cycles between classes

You probably need to split a class or two.
Often the cycles appear in classes [which perform several tasks](https://refactoring.guru/smells/large-class).

## Settings

### --depth
Number of steps to go further down into the package hierarchy starting at the root package.
The root package is detected automatically.

### --include

Regular expression to filter files by their full path before analysis.
All paths containing a match are analyzed.
You can use this to exclude auto-generated clients.

```shell
sda visualize --include 'src/main/java'
```

### --max-cycles (validate only)
Maximum number of cycles to attribute with exit-code _0_.

The parameter '--max-cycles' is intended as follows:
* remove cycles step-by-step from legacy projects with the goal to set --max-cycles to zero eventually
* rare corner-cases where you consider cycles a valid option

### --graphLabel (visualize only)
The graph label is located at the bottom center of the resulting image (default "rendered by github.com/sandstorm/dependency-analysis").

### --output (visualize only)
Path to the image file to generate, use `stdout` to output DOT graph without image rendering.

### --show-image (visualize and OSX only)
Automatically open the image after rendering.
Internally `open …` is called.
On operating systems other then OSX this flag and feature does not exist.

```shell
sda visualize --show-image=false
```

### --show-node-labels (visualize only)
Render graph with node labels.
I use this flag to generate anonymous graphs for this documentation.

###  -T, --type string (visualize only)
Type of the image file.

```
$ sda visualize listSupportedOutputTypes
bmp       BMP Windows Bitmap
cgimage   CGImage Apple Core Graphics
canon     DOT Graphviz Language
dot       DOT Graphviz Language
gv        DOT Graphviz Language
xdot      DOT Graphviz Language
xdot1.2   DOT Graphviz Language
xdot1.4   DOT Graphviz Language
eps       EPS Encapsulated PostScript
exr       EXR OpenEXR
fig       FIG Xfig
gd        GD LibGD
gd2       GD2 LibGD
gif       GIF Graphics Interchange Format
gtk       GTK Formerly GTK+ / GIMP ToolKit
ico       ICO Windows Icon
cmap      Image Map: Client-side
ismap     Image Map: Server-side
imap      Image Map: Server-side and client-side
cmapx     Image Map: Server-side and client-side
imap_np   Image Map: Server-side and client-side
cmapx_np  Image Map: Server-side and client-side
jpg       JPEG Joint Photographic Experts Group
jpeg      JPEG Joint Photographic Experts Group
jpe       JPEG Joint Photographic Experts Group
jp2       JPEG 2000
json      JSON JavaScript Object Notation
json0     JSON JavaScript Object Notation
dot_json  JSON JavaScript Object Notation
xdot_json JSON JavaScript Object Notation
pdf       PDF Portable Document Format
pic       PIC Brian Kernighan's Diagram Language
pct       PICT Apple PICT
pict      PICT Apple PICT
plain     Plain Text Simple, line-based language
plain-ext Plain Text Simple, line-based language
png       PNG Portable Network Graphics
pov       POV-Ray Persistence of Vision Raytracer (prototype)
ps        PS Adobe PostScript
ps2       PS/PDF Adobe PostScript for Portable Document Format
psd       PSD Photoshop
sgi       SGI Silicon Graphics Image
svg       SVG Scalable Vector Graphics
svgz      SVG Scalable Vector Graphics
tga       TGA Truevision TARGA
tif       TIFF Tag Image File Format
tiff      TIFF Tag Image File Format
tk        Tk Tcl/Tk
vm        VML Vector Markup Language
vmlz      VML Vector Markup Language
vrml      VRML Virtual Reality Modeling Language
wbmp      WBMP Wireless Bitmap
webp      WebP WebP
xlib      X11 X11 Window
x11       X11 X11 Window
````

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

# Developer documentation

* [install Golang](https://golangdocs.com/install-go-mac-os)
* `brew install graphviz`
* [install the Sandstorm `dev` script helper](https://github.com/sandstorm/dev-script-runner)

## Start from source

```sh
dev run
# or
dev run --help
```

## Release

```sh
# adjust version as needed
git tag v1.0.0
dev release
```

## Add a new command

```sh
# install COBRA CLI (see links to docs below)
cobra-cli add helloWorld
# and adjust generated files
```

## Helpful development commands

```sh
# see all dev scripts
dev help

# visualize all Java projects in current folder
find . -type d | grep -iE 'src/main/java$' | xargs -Isrc bash -c 'export project=$(echo src | cut -d / -f 2); sda visualize src -o $project.svg -l $project'
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
