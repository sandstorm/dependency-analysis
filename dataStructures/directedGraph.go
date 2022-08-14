package dataStructures

// directed graph with nodes of type string
type DirectedStringGraph struct {
	// all nodes with their set of children (might be empty)
	Edges map[string]*StringSet
}

func NewDirectedStringGraph() *DirectedStringGraph {
	result := &DirectedStringGraph{}
	result.Edges = make(map[string]*StringSet)
	return result
}

// Adds a node to the graph unless it already exists
func (this *DirectedStringGraph) AddNode(node string) *DirectedStringGraph {
	if this.Edges[node] == nil {
		this.Edges[node] = NewStringSet()
	}
	return this
}

// Adds an edge to the graph unless it already exists
func (this *DirectedStringGraph) AddEdge(start string, target string) *DirectedStringGraph {
	this.AddNode(start)
	this.AddNode(target)
	this.Edges[start].Add(target)
	return this
}

// Provides all children in a set or nil for nodes not in this graph
func (this *DirectedStringGraph) GetChildren(node string) *StringSet {
	return this.Edges[node]
}

func (this *DirectedStringGraph) GetNodes() []string {
	result := make([]string, len(this.Edges))
	var i = 0
	for node, _ := range this.Edges {
		result[i] = node
		i++
	}
	return result
}

func (this *DirectedStringGraph) GetEdges() map[string][]string {
	result := make(map[string][]string, len(this.Edges))
	for caller, callees := range this.Edges {
		result[caller] = callees.ToArray()
	}
	return result
}
