package disjointset

// DisjointSet is the interface for the disjoint-set (or union-find) data
// structure.
// Do not change the definition of this interface.
type DisjointSet interface {
	// UnionSet(s, t) merges (unions) the sets containing s and t,
	// and returns the representative of the resulting merged set.
	UnionSet(int, int) int
	// FindSet(s) returns representative of the class that s belongs to.
	FindSet(int) int
}

type disjointSet struct {
	parent map[int]int
	size   map[int]int
}

func NewDisjointSet() DisjointSet {
	return &disjointSet{
		parent: make(map[int]int),
		size:   make(map[int]int),
	}
}

func (d *disjointSet) FindSet(u int) int {
	if _, ok := d.parent[u]; !ok {
		d.parent[u] = u
		d.size[u] = 1
		return u
	}
	// Path compression
	for d.parent[u] != u {
		d.parent[u] = d.parent[d.parent[u]]  // Path compression step
		u = d.parent[u]
	}
	return u
}

func (d *disjointSet) UnionSet(a, b int) int {
	rootA := d.FindSet(a)
	rootB := d.FindSet(b)

	if rootA == rootB {
		return rootA
	}

	// Union by size
	if d.size[rootA] < d.size[rootB] {
		d.parent[rootA] = rootB
		d.size[rootB] += d.size[rootA]
		return rootB
	} else {
		d.parent[rootB] = rootA
		d.size[rootA] += d.size[rootB]
		return rootA
	}
}