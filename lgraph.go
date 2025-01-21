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
func FindSequence(g1, g2 LGraph, s, t node, k uint) ([]rune, bool) {
	// TODO: Complete the function.
	panic("TODO: implement this!")
	pathG1 := findPaths(g1, s, t, k) // find all paths in g1 from s to t with len k
	pathG2 := findPaths(g2, s, t, k) // find all paths in g2 from s to t with len k

	// if there is a path in G1 also in G2 with the same len k, then return false
	for _, path1 := range pathG1 {
		for _, path2 := range pathG2 {
			if path1 == path2 {
				return nil, false
			}
		}
	}
	// otherwise, return true
	return []rune(pathG1[0]), true
}

// helper function
func findPaths(g LGraph, s, t node, k uint) []string {
	var paths []string // a slice that store all path
	// using dfs to find the path from s to t
	var dfs func(current node, len uint, path string)
	dfs = func(current node, len uint, path string) {
		// base case
		if len > k {
			return
		}
		// reach the last node
		if len == k && current == t {
			paths = append(paths, path)
			return
		}

		edges, exists := g(current)
		if !exists {
			return
		}

		for _, edge := range edges {
			dfs(edge.destination, len+1, path+string(edge.label))
		}
	}
	dfs(s, 0, "")
	return paths
}

/*
	Generat all path in G1
	- start from node s and find all path of len k that end at node t
	- store each path as a string (eg "abc") in a slice for G1

	Generate all path in G2
	- similarly to generate path in G1, start from s to t with the len k
	- store each path as string in a slice for g2
	Compare paths from g1 and g2
	- for each path in g1, if there is a path that exists in g2 -> false
	- otherwise, return the path and true

	build paths -> compare paths -> return result
*/
