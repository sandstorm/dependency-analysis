package analyser

import (
	"sort"
	"github.com/sandstorm/dependency-analysis/utils"
)

// directed graph with nodes of type string
// nodes may contain weights
type WeightedStringGraph struct {
	DirectedStringGraph
	// weights per node
	weights map[string]int
}

func NewWeightedStringGraph() *WeightedStringGraph {
    result := &WeightedStringGraph{}
    result.edges = make(map[string]*utils.StringSet)
    result.weights = make(map[string]int)
    return result
}

func (this *WeightedStringGraph) SetWeight(node string, weight int) {
	this.weights[node] = weight
}

func (this *WeightedStringGraph) GetWeight(node string) int {
	return this.weights[node]
}

func (this *WeightedStringGraph) HasWeight(node string) bool {
	_, hasWeight := this.weights[node]
	return hasWeight
}

// Creates an new weightes graph with the same structure as the
// given graph. Weights are the number grand*-children of each node.
// Cycles are tolerated and each node at most counts as one grand*-child.
func WeightByNumberOfDescendant(source *DirectedStringGraph) *WeightedStringGraph {
	result := NewWeightedStringGraph()
	result.edges = source.edges
	nodes := result.getNodes()
	// If the graph contains cycles the result depends on the order of the nodes
	// during the iteration. The order in the set is not defined.
	// So we sort the nodes here to have a deterministic (though kind of arbitrary)
	// result.
	sort.Strings(nodes)
	for _, node := range nodes {
		if !result.HasWeight(node) {
			result.calculateWeightsByDescendants(node, utils.NewStringSet())
		}
	}
	return result
}

// Recursively sets the weights of the given node and all its descendants tto the number
// of reachable distinct nodes.
func (this *WeightedStringGraph) calculateWeightsByDescendants(node string, visitedNodes *utils.StringSet) int {
	if this.HasWeight(node) {
		return this.GetWeight(node)
	}
	visitedNodes.Add(node)
	// we must remove this node from the list of visited nodes when
	// the function returns since it then has a weight already
	// and must not be skipped any more
	defer visitedNodes.Remove(node)

	allChildren := this.GetChildren(node).ToArray()
	var weight = 0
	for _, child := range allChildren {
		if !visitedNodes.Contains(child) {
			maybeWeight := 1 + this.calculateWeightsByDescendants(child, visitedNodes)
			if maybeWeight > weight {
				weight = maybeWeight
			}
		}
	}
	this.SetWeight(node, weight)
	return weight
}
