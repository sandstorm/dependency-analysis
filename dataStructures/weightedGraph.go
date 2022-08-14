package dataStructures

// directed graph with nodes of type string
// nodes may contain weights
type WeightedStringGraph struct {
	DirectedStringGraph
	weightsByNode map[string]int
}

func NewWeightedStringGraph() *WeightedStringGraph {
	result := &WeightedStringGraph{}
	result.Edges = make(map[string]*StringSet)
	result.weightsByNode = make(map[string]int)
	return result
}

func (this *WeightedStringGraph) SetWeight(node string, weight int) {
	this.weightsByNode[node] = weight
}

func (this *WeightedStringGraph) GetWeight(node string) int {
	return this.weightsByNode[node]
}

func (this *WeightedStringGraph) HasWeight(node string) bool {
	_, hasWeight := this.weightsByNode[node]
	return hasWeight
}

func (this *WeightedStringGraph) GetNodesGroupedByWeight() (map[int][]string, int) {
	nodesByWeights := make(map[int]*StringSet)
	for _, node := range this.GetNodes() {
		weight := this.GetWeight(node)
		var nodes, isSet = nodesByWeights[weight]
		if !isSet {
			nodes = NewStringSet()
			nodesByWeights[weight] = nodes
		}
		nodes.Add(node)
	}
	maxWeight := 0
	result := make(map[int][]string)
	for weight, nodes := range nodesByWeights {
		result[weight] = nodes.ToArray()
		if weight > maxWeight {
			maxWeight = weight
		}
	}
	return result, maxWeight
}
