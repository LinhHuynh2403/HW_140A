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

// TODO: implement a type that satisfies the DisjointSet interface.

// DisjointSetImpl satisfies the DisjointSet interface.
// A node is a representative if it points to itself in 'parent map.
// An undefined key in 'parent' should be taken as pointing to itself.
// The size of a class is stored in the representative's 'size map.
// https://en.wikipedia.org/wiki/Disjoint-set_data_structure
type DisjointSetImpl struct {
	parent map[int]int
	size   map[int]int
}

// UnionSet merges the classes represented by s and t, using Union by size, I
// and returns the new class representative.
func (ds *DisjointSetImpl) UnionSet(t, s int) int {
	// We begin by calling FindSet() on the two integers representing the classes we want to union.
	// This allows us to perform union operations using their class representatives.
	s, t = ds.FindSet(s), ds.FindSet(t)
	// If we were trying to union two equivalent sets, then their class representatives would be the same, and we
	// can just return one of them.
	if s == t {
		return s
	}
	// Save the respective sizes of the two sets so we know which one is bigger.
	sizeS, sizeT := ds.size[s], ds.size[t]
	// If sizeT is bigger, do a swap since we are going to be making s the new representative.
	if sizeS < sizeT {
		s, t = t, s
	}
	// Now the parent of t is s and the size of s is its original size plus the size of t's set.
	ds.parent[t] = s
	ds.size[s] = sizeS + sizeT
	return s
}

// FindSet returns representative of the class that s belongs to.
func (ds *DisjointSetImpl) FindSet(s int) int {
	// Compress path for all non-representative nodes visited.

	// If the int s does not yet exist in our disjoint set parent map, it is assumed to be a singleton set.
	// We initialize it so that its parent is itself and its size is 1.
	if p, ok := ds.parent[s]; !ok {
		ds.parent[s] = s
		ds.size[s] = 1
		return s

		// If the parent of s is not itself, then we call FindSet() on the parent, reassigning the parent in the map to the root.
		// This is the heart of the path compression optimization.
	} else if s != p {
		r := ds.FindSet(p)
		ds.parent[s] = r
		return r
		// Otherwise, we just return s.
	} else {
		return s
	}
}

// NewDisjointSet creates a struct of a type that satisfies the DisjointSet interface.
func NewDisjointSet() DisjointSet {
	// panic("TODO: implement this!")
	return &DisjointSetImpl{
		parent: make(map[int]int),
		size:   make(map[int]int),
	}
}
