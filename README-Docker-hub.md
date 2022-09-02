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
