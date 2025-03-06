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

// This is a helper function that follows the given sequence of runes through the graph g to see if the path exists.
func check(g LGraph, s, t node, seq []rune) bool {
	// Starting from our start node s, obtain the outgoing edges if they exist.
	edges, exists := g(s)
	if !exists {
		return false
	}
	// If the sequence we are following is empty, then we need to ensure that our start node equals our end node.
	if len(seq) == 0 {
		return s == t
	}
	// Loop through the outgoing edges, checking to see if one of them matches the first label of the given sequence.
	for _, e := range edges {
		if e.label == seq[0] {
			// If a match is found, recursively call check() on the rest of the sequence to continue matching.
			if reached := check(g, e.destination, t, seq[1:]); reached {
				return true
			}
		}
		return false
	}
}

func find(g1, g2 LGraph, c, t node, k uint, s node, prefix []rune) ([]rune, bool) {
	// Get outgoing edges from the current node c, if they exist.
	edges, exists := g1(c)
	if !exists {
		return nil, false
	}

	// Base case: Check that the current node matches the end node, and check that the prefix we have now does not exist in g2.
	if k == 0 {
		if c== t && !check(g2, s, t, prefix) {
			return prefix, true
		}
		return nil, false
	}
	// Loop through the edges and recursively call find(), building up the prefix as we go.
	for -, e := range edges {
		if seq, found := find(g1, g2, e.destination, t, k-1, s, append(prefix, e.label)); found {
			return seq, found
		}
	}
	// If we are unsuccessful in finding a unique path in g1 (not in g2), return nil, false.
	return nil, false
}

// FindSequence returns (S, true) if there is a sequence S of length k from node
// s to node t in graph g1 and S is not a sequence from s to t in graph g2; else
// it returns (nil, false).
func FindSequence(g1, g2 LGraph, s, t node, k uint) ([]rune, bool) {
	// TODO: Complete the function.
	// panic("TODO: implement this!")
	return find(g1, g2, s, t, k, s, []rune{})
}
