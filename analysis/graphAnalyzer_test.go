package analysis

import (
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, message string, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(actual, expected) {
		return
	}
	t.Errorf("%s\nexpected %v (type %v), received %v (type %v)",
		message,
		expected, reflect.TypeOf(expected),
		actual, reflect.TypeOf(actual))
}

func TestWeightByNumberOfDescendant(t *testing.T) {
	testCases := []struct {
		name     string
		graph    *dataStructures.DirectedStringGraph
		expected map[string]int
	}{
		{
			name:     "empty graph",
			graph:    dataStructures.NewDirectedStringGraph(),
			expected: map[string]int{},
		},
		{
			name: `
			| ┌─────┐
			| │  A  │
			| │  0  │
			| └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddNode("A"),
			expected: map[string]int{
				"A": 0,
			},
		},
		{
			name: `
			| ┌─────┐    ┌─────┐
			| │  A  │    │  B  │
			| │  0  │    │  0  │
			| └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddNode("A").
				AddNode("B"),
			expected: map[string]int{
				"A": 0,
				"B": 0,
			},
		},
		{
			name: `
			| ┌─────┐    ┌─────┐
			| │  A  │    │  B  │
			| │  1  │───▶│  0  │
			| └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B"),
			expected: map[string]int{
				"A": 1,
				"B": 0,
			},
		},
		{
			name: `
            | ┌─────┐    ┌─────┐
            | │  A  │    │  B  │
            | │  1  │───▶│  0  │
            | └─────┘    └─────┘
            | ┌─────┐    ┌─────┐
            | │  a  │    │  b  │
            | │  1  │───▶│  0  │
            | └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("a", "b"),
			expected: map[string]int{
				"A": 1,
				"B": 0,
				"a": 1,
				"b": 0,
			},
		},
		{
			name: `
            | ┌─────┐      ┌─────┐       
            | │  A  │      │  B  │       
            | │  3  │─────▶│  2  │       
            | └─────┘      └─────┘       
            |                │           
            |                │          
            |      ┌─────┐   │   ┌─────┐
            |      │  C  │   │   │  D  │
            |      │  1  │◀──┴──▶│  1  │
            |      └─────┘       └─────┘
            |          │             │   
            |          │             ▼   
            | ┌─────┐  │  ┌─────┐ ┌─────┐
            | │  E  │  │  │  F  │ │  G  │
            | │  0  │◀─┴─▶│  0  │ │  0  │
            | └─────┘     └─────┘ └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "C").
				AddEdge("B", "D").
				AddEdge("C", "E").
				AddEdge("C", "F").
				AddEdge("D", "G"),
			expected: map[string]int{
				"A": 3,
				"B": 2,
				"C": 1,
				"D": 1,
				"E": 0,
				"F": 0,
				"G": 0,
			},
		},
		{
			name: `
			| ┌─────┐       ┌─────┐          
			| │  A  │       │  B  │          
			| │  3  │──────▶│  2  │────────┐
			| └─────┘       └─────┘        │
			|                 │            │
			|                 │            │
			|       ┌─────┐   │   ┌─────┐  │
			|       │  D  │   │   │  C  │  │
			|       │  1  │◀──┴──▶│  1  │  │
			|       └─────┘       └─────┘  │
			|          │             │     │
			|          │             ▼     │
			| ┌─────┐  │  ┌─────┐ ┌─────┐  │
			| │  E  │  │  │  F  │ │  G  │  │
			| │  0  │◀─┴─▶│  0  │ │  0  │◀─┘
			| └─────┘     └─────┘ └─────┘   
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "C").
				AddEdge("B", "D").
				AddEdge("B", "G").
				AddEdge("C", "E").
				AddEdge("C", "F").
				AddEdge("D", "G"),
			expected: map[string]int{
				"A": 3,
				"B": 2,
				"C": 1,
				"D": 1,
				"E": 0,
				"F": 0,
				"G": 0,
			},
		},
		{
			name: `
    		|    ┌────────────┐   
			|    │            │   
			|    ▼            │   
			| ┌─────┐      ┌─────┐
			| │  A  │      │  B  │
			| │  1  │      │  0  │
			| └─────┘      └─────┘
			|    │            ▲   
			|    │            │   
			|    └────────────┘   
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A"),
			expected: map[string]int{
				"A": 1,
				"B": 0,
			},
		},
		{
			name: `
            | ┌─────┐      ┌─────┐         
            | │  A  │      │  B  │         
            | │  3  │─────▶│  2  │◀────┐   
            | └─────┘      └─────┘     │   
            |                 │        │   
            |                 │        │   
            |                 ▼        │   
            |              ┌─────┐     │   
            |              │  C  │     │   
            |              │  1  │     │   
            |              └─────┘     │   
            |                │         │   
            |     ┌─────┐    │     ┌─────┐
            |     │  D  │    │     │  E  │
            |     │  0  │◀───┴────▶│  0  │
            |     └─────┘          └─────┘
            `,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "C").
				AddEdge("C", "D").
				AddEdge("C", "E").
				AddEdge("E", "B"),
			expected: map[string]int{
				"A": 3,
				"B": 2,
				"C": 1,
				"D": 0,
				"E": 0,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run("testing WeightByNumberOfDescendant…", func(t *testing.T) {
			t.Log(testCase.name)
			actualGraph := WeightByNumberOfDescendant(testCase.graph)
			for node, expected := range testCase.expected {
				AssertEquals(t,
					"incorrect weight of node "+node,
					expected,
					actualGraph.GetWeight(node))
			}
		})
	}
}

func TestFindCycles(t *testing.T) {
	testCases := []struct {
		name     string
		graph    *dataStructures.DirectedStringGraph
		expected []dataStructures.Cycle
	}{
		{
			name:     "empty graph",
			graph:    dataStructures.NewDirectedStringGraph(),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
			| ┌─────┐
			| │  A  │
			| └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddNode("A"),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
			| ┌─────┐    ┌─────┐
			| │  A  │    │  B  │
			| └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddNode("A").
				AddNode("B"),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
			| ┌─────┐    ┌─────┐
			| │  A  │───▶│  B  │
			| └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B"),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
            | ┌─────┐    ┌─────┐
            | │  A  │───▶│  B  │
            | └─────┘    └─────┘
            | ┌─────┐    ┌─────┐
            | │  a  │───▶│  b  │
            | └─────┘    └─────┘
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("a", "b"),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
    		|    ┌────────────┐   
			|    ▼            │   
			| ┌─────┐      ┌─────┐
			| │  A  │      │  B  │
			| └─────┘      └─────┘
			|    │            ▲   
			|    └────────────┘   
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"A": "B",
					"B": "A",
				},
			},
		},
		{
			name: `
			|   ┌───────┐  
			|   ▼       │  
			| ┌───┐   ┌───┐
			| │ A │   │ B │
			| └───┘   └───┘
			|   │       ▲  
			|   └───────┘  
			|   ┌───────┐  
			|   ▼       │  
			| ┌───┐   ┌───┐
			| │ a │   │ b │
			| └───┘   └───┘
			|   │       ▲  
			|   └───────┘  
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A").
				AddEdge("a", "b").
				AddEdge("b", "a"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"A": "B",
					"B": "A",
				},
				map[string]string{
					"a": "b",
					"b": "a",
				},
			},
		},
		{
			name: `
			|   ┌────┐
			|   ▼    │
			| ┌───┐  │
			| │ A │──┘
			| └───┘   
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "A"),
			expected: []dataStructures.Cycle{},
		},
		{
			name: `
			|   ┌───────┐        ┌───┐
			|   ▼       │    ┌───│ C │
			| ┌───┐   ┌───┐  │   └───┘
			| │ A │   │ B │◀─┤        
			| └───┘   └───┘  │   ┌───┐
			|   │       ▲    └───│ D │
			|   └───────┘        └───┘          
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A").
				AddEdge("C", "B").
				AddEdge("D", "B"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"A": "B",
					"B": "A",
				},
			},
		},
		{
			name: `
			|           ┌──────────┐  
			|           │          │  
			|   ┌───────┤          ▼  
			|   │       │        ┌───┐
			|   ▼       │    ┌───│ C │
			| ┌───┐   ┌───┐  │   └───┘
			| │ A │   │ B │◀─┤        
			| └───┘   └───┘  │   ┌───┐
			|   │       ▲    └───│ D │
			|   │       │        └───┘
			|   └───────┘             
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A").
				AddEdge("B", "C").
				AddEdge("C", "B").
				AddEdge("D", "B"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"A": "B",
					"B": "A",
				},
				map[string]string{
					"B": "C",
					"C": "B",
				},
			},
		},
		{
			name: `
			|   ┌───────┬────────┐  
			|   │       │        │  
			|   ▼       │        ▼  
			| ┌───┐   ┌───┐    ┌───┐
			| │ A │   │ B │    │ C │
			| └───┘   └───┘    └───┘
			|   │       ▲        │  
			|   │       │        │  
			|   └───────┴────────┘  
			`,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "A").
				AddEdge("B", "C").
				AddEdge("C", "B"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"A": "B",
					"B": "A",
				},
				map[string]string{
					"B": "C",
					"C": "B",
				},
			},
		},
		{
			name: `
            | ┌─────┐      ┌─────┐         
            | │  A  │─────▶│  B  │◀────┐   
            | └─────┘      └─────┘     │   
            |                 │        │   
            |                 ▼        │   
            |              ┌─────┐     │   
            |              │  C  │     │   
            |              └─────┘     │   
            |     ┌─────┐    │     ┌─────┐
            |     │  D  │◀───┴────▶│  E  │
            |     └─────┘          └─────┘
            `,
			graph: dataStructures.NewDirectedStringGraph().
				AddEdge("A", "B").
				AddEdge("B", "C").
				AddEdge("C", "D").
				AddEdge("C", "E").
				AddEdge("E", "B"),
			expected: []dataStructures.Cycle{
				map[string]string{
					"B": "C",
					"C": "E",
					"E": "B",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run("testing FindCycles", func(t *testing.T) {
			t.Log(testCase.name)
			actualCycles := FindCycles(testCase.graph)
			for i := 0; i < min(len(actualCycles), len(testCase.expected)); i++ {
				AssertEquals(t,
					"incorrect cycle",
					testCase.expected[i],
					actualCycles[i])
			}
			if len(testCase.expected) != len(actualCycles) {
				t.Errorf("incorrect number of cycles\nexpected %d, received %d",
					len(testCase.expected),
					len(actualCycles))
			}
		})
	}
}
