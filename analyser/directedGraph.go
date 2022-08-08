package analyser

import (
	"github.com/sandstorm/dependency-analysis/utils"
)

// directed graph with nodes of type string
type DirectedStringGraph struct {
	// all nodes with their set of children (might be empty)
    edges map[string]*utils.StringSet
}

func NewDirectedStringGraph() *DirectedStringGraph {
    result := &DirectedStringGraph{}
    result.edges = make(map[string]*utils.StringSet)
    return result
}

// Adds a node to the graph unless it already exists
func (this *DirectedStringGraph) AddNode(node string) *DirectedStringGraph {
	if this.edges[node] == nil {
		this.edges[node] = utils.NewStringSet()
	}
	return this
}

// Adds an edge to the graph unless it already exists
func (this *DirectedStringGraph) AddEdge(start string, target string) *DirectedStringGraph {
	this.AddNode(start)
	this.AddNode(target)
	this.edges[start].Add(target)
	return this
}

// Provides all children in a set or nil for nodes not in this graph
func (this *DirectedStringGraph) GetChildren(node string) *utils.StringSet {
	return this.edges[node]
}

func (this *DirectedStringGraph) getNodes() []string {
	result := make([]string, len(this.edges))
	var i = 0
	for node, _ := range this.edges {
		result[i] = node
		i++
	}
	return result
}