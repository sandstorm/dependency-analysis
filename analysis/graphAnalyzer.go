package analysis

import (
	"sort"
	"container/list"
	"github.com/sandstorm/dependency-analysis/dataStructures"
)

// Creates an new weightes graph with the same structure as the
// given graph. Weights are the number grand*-children of each node.
// Cycles are tolerated and each node at most counts as one grand*-child.
func WeightByNumberOfDescendant(source *dataStructures.DirectedStringGraph) *dataStructures.WeightedStringGraph {
	result := dataStructures.NewWeightedStringGraph()
	result.Edges = source.Edges
	nodes := result.GetNodes()
	// If the graph contains cycles the result depends on the order of the nodes
	// during the iteration. The order in the set is not defined.
	// So we sort the nodes here to have a deterministic (though kind of arbitrary)
	// result.
	sort.Strings(nodes)
	for _, node := range nodes {
		if !result.HasWeight(node) {
			calculateWeightsByDescendants(result, node, dataStructures.NewStringSet())
		}
	}
	return result
}

// Recursively sets the weights of the given node and all its descendants tto the number
// of reachable distinct nodes.
func calculateWeightsByDescendants(target *dataStructures.WeightedStringGraph, node string, visitedNodes *dataStructures.StringSet) int {
	if target.HasWeight(node) {
		return target.GetWeight(node)
	}
	visitedNodes.Add(node)
	// we must remove this node from the list of visited nodes when
	// the function returns since it then has a weight already
	// and must not be skipped any more
	defer visitedNodes.Remove(node)

	allChildren := target.GetChildren(node).ToArray()
	var weight = 0
	for _, child := range allChildren {
		if !visitedNodes.Contains(child) {
			maybeWeight := 1 + calculateWeightsByDescendants(target, child, visitedNodes)
			if maybeWeight > weight {
				weight = maybeWeight
			}
		}
	}
	target.SetWeight(node, weight)
	return weight
}

// Provides all cyclic paths.
// 
// Runs a DFS (depth first search) over all path in the graph capturing
// found cycles. All cycles are collected with two exception:
// 
// 1) Cycles with only one node
//
// We analyze dependencies here and dependencies within one package are ok.
// 
// 2) Cycles using the same node several times
//
// Cycles leaving the same node more than once are not considered, example:
// A graph in the shape of an eight contains three cycles, but we are
// only interested in the two shorter ones.
//   ┌───────┬────────┐  
//   │       │        │  
//   ▼       │        ▼  
// ┌───┐   ┌───┐    ┌───┐
// │ A │   │ B │    │ C │
// └───┘   └───┘    └───┘
//   │       ▲        │  
//   │       │        │  
//   └───────┴────────┘  
// Cycles:
// * A -> B -> A (considered)
// * C -> B -> C (considered)
// * A -> B -> C -> B -> A (ignored)
func FindCycles(graph *dataStructures.DirectedStringGraph) []dataStructures.Cycle {
	nodes := graph.GetNodes()
	// We want deterministic results but the result
	// depends on the iteration order of the graph
	// which is undefined.
	sort.Strings(nodes)
	currentPath := dataStructures.NewStringArrayList(len(nodes))
	doneNodes := dataStructures.NewStringSet()
	resultList := list.New()
	for _, node := range nodes {
		executeFindCycles(graph, node, currentPath, doneNodes, resultList)
	}

	result := make([]dataStructures.Cycle, resultList.Len())
	var i = 0
	for element := resultList.Front(); element != nil; element = element.Next() {
		if value, ok := element.Value.(dataStructures.Cycle); ok {
			result[i] = value
			i++
		} 
    }
	return result
}

// internal implementation of FindCycles
func executeFindCycles(
	graph *dataStructures.DirectedStringGraph,
	currentNode string,
	currentPath *dataStructures.StringArrayList,
	doneNodes *dataStructures.StringSet,
	result *list.List) {
	// already searched ?
	if doneNodes.Contains(currentNode) {
		// If we reach a node which already has been analyzed we can abort the search,
		// because:
		// We have analyzed the current node and already found all cycles it is part of.
		// This already includes any cycle starting with "currentPath", thus we can stop now.
		return
	}
	// found cycle ? 
	indexOnCurrentPath := currentPath.IndexOf(currentNode)
	if indexOnCurrentPath >= 0 {
		// found cycle
		cycle := make(dataStructures.Cycle)
		for i := indexOnCurrentPath; i < currentPath.Length(); i++ {
			next := i + 1
			if next < currentPath.Length() {
				cycle[currentPath.Get(i)] = currentPath.Get(next)
			} else {
				cycle[currentPath.Get(i)] = currentNode
			}
		}
		result.PushBack(cycle)
		return
	}
	// proceed recursion
	currentPath.Push(currentNode)
	defer currentPath.Pop(currentNode)
	defer doneNodes.Add(currentNode)
	children := graph.GetChildren(currentNode).ToArray()
	// We want deterministic results but the result
	// depends on the iteration order of the graph
	// which is undefined.
	sort.Strings(children)
	for _, child := range children {
		if child != currentNode {
			executeFindCycles(graph, child, currentPath, doneNodes, result)
		}
	}
}
