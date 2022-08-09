/*
This package renders graphs into the DOT language of GraphViz.

For rendering we use GraphViz so we don't have to/cannot control the position of
each and every node. However we create groups of nodes to improve the overall
layout.

The GraphViz CLI tools render the final SVG/PNG/… so please execute
$ brew install graphviz
or see https://graphviz.org/download/
*/
package rendering
