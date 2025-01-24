package lgraph

type node uint

type edge struct {
	destination node
	label       rune
}

// LGraph is a function representing a directed labeled graph. If the node exists
// in the graph, the function returns true along with the set of outgoing edges
// from that node, otherwise false and nil.
type LGraph func(node) ([]edge, bool)

// FindSequence returns (S, true) if there is a sequence S of length k from node
// s to node t in graph g1 and S is not a sequence from s to t in graph g2; else
// it returns (nil, false).
func FindSequence(g1, g2 LGraph, s node, t node, k uint) ([]rune, bool) {
	sequences := generateSequences(g1, s, t, k)
	for _, seq := range sequences {
		if !isSequencePresent(g2, s, t, seq) {
			return seq, true
		}
	}
	return nil, false
}

// helper function
func generateSequences(g LGraph, s node, t node, k uint) [][]rune {
	var result [][]rune
	var dfs func(current node, path []rune, steps uint)

	dfs = func(current node, path []rune, steps uint) {
		if steps == k {
			if current == t {
				newPath := make([]rune, len(path))
				copy(newPath, path)
				result = append(result, newPath)
			}
			return
		}

		edges, exists := g(current)
		if !exists {
			return
		}

		
		for _, e := range edges {
			dfs(e.destination, append(path, rune(e.label)), steps+1)
		}
	}

	dfs(s, []rune{}, 0)
	return result
}