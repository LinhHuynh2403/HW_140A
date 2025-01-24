package lgraph

type node uint

type edge struct {
	destination node
	label       rune
}

type LGraph func(node) ([]edge, bool)

func FindSequence(g1, g2 LGraph, s node, t node, k uint) ([]rune, bool) {
	sequences := generateSequences(g1, s, t, k)
	for _, seq := range sequences {
		if !isSequencePresent(g2, s, t, seq) {
			return seq, true
		}
	}
	return nil, false
}

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

func isSequencePresent(g LGraph, s node, t node, sequence []rune) bool {
    if len(sequence) == 0 {
        _, exists := g(s)
        return exists && s == t
    }

    currentNodes := make(map[node]bool)
    currentNodes[s] = true

    for _, label := range sequence {
        nextNodes := make(map[node]bool)

        for n := range currentNodes {
            edges, exists := g(n)
            if !exists {
                continue
            }
            for _, e := range edges {
                if e.label == label {
                    nextNodes[e.destination] = true
                }
            }
        }

        if len(nextNodes) == 0 {
            return false
        }
        currentNodes = nextNodes
    }

    return currentNodes[t]
}