package lgraph

import (
    "sync"
)

type node uint

type edge struct {
    destination node
    label       rune
}

type LGraph func(node) ([]edge, bool)

func FindSequence(g1, g2 LGraph, s node, t node, k uint) ([]rune, bool) {
    sequences := generateSequencesConcurrent(g1, s, t, k)

    resultChan := make(chan []rune, len(sequences))
    var wg sync.WaitGroup

    for _, seq := range sequences {
        wg.Add(1)
        go func(seq []rune) {
            defer wg.Done()
            if !isSequencePresent(g2, s, t, seq) {
                resultChan <- seq
            }
        }(seq)
    }

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    for seq := range resultChan {
        return seq, true
    }

    return nil, false
}

func generateSequencesConcurrent(g LGraph, s node, t node, k uint) [][]rune {
    if k == 0 {
        _, exists := g(s)
        if exists && s == t {
            return [][]rune{{}}
        }
    }

    var result [][]rune
    resultChan := make(chan []rune, 100)
    var wg sync.WaitGroup
    var mu sync.Mutex

    var dfs func(current node, path []rune, steps uint)
    dfs = func(current node, path []rune, steps uint) {
        if steps == k {
            if current == t {
                newPath := make([]rune, len(path))
                copy(newPath, path)
                resultChan <- newPath
            }
            return
        }

        edges, exists := g(current)
        if !exists {
            return
        }

        for _, e := range edges {
            wg.Add(1)
            go func(e edge, pathCopy []rune) {
                defer wg.Done()
                newPath := append([]rune{}, pathCopy...)
                newPath = append(newPath, e.label)
                dfs(e.destination, newPath, steps+1)
            }(e, path)
        }
    }

    wg.Add(1)
    go func() {
        defer wg.Done()
        dfs(s, []rune{}, 0)
    }()

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    for seq := range resultChan {
        mu.Lock()
        result = append(result, seq)
        mu.Unlock()
    }

    return result
}

func isSequencePresent(g LGraph, s node, t node, seq []rune) bool {
    if len(seq) == 0 {
        _, exists := g(s)
        return exists && s == t
    }

    current := s
    for _, label := range seq {
        edges, exists := g(current)
        if !exists {
            return false
        }

        found := false
        for _, e := range edges {
            if e.label == label {
                current = e.destination
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return current == t
}