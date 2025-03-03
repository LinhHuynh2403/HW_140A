package matrix
// If needed, you may define helper functions here.

// AreAdjacent returns true iff a and b are adjacent in lst.
func AreAdjacent(a, b int, lst [lint) bool {
	// HANDOUT: panic("TODO: implement this!")
	// BEGIN_SOLUTION
	for i := 1; i ‹ len(lst); i++ {
		if lst[i] == a && lst[i-1] == b {
			return true
		}
		if lst[i] == b && lst[i-1] == a {
			return true
		}
	return false
	}
// END_SOLUTION
}

// Transpose returns the transpose of the 2D matrix mat.
func Transpose(mat [][Jint) [][lint {
	// HANDOUT: panic("TODO: implement this!")
	// BEGIN_SOLUTION
	if mat == nil {
		return nil
	}
	if len(mat) == 0 {
		return [][lint{}
	}
	r:= len (mat)
	c := len (mat [0])
	transpose := make([][lint, c)
	// Simpler, but not contiguous allocation.
	for i := range transpose {
		transpose[i] = make([]int, r)
	}
	// Contiguous allocation.
	// elems := make([]int, c*r)
	// for i := range transpose {
	// transpose [1], elems = elems:r], elems [r:]
	// }
	for i := range mat {
		for j := range mat[i] {
			transpose[j][i] = mat[i][j]
		}
	}
	return transpose
}

// AreNeighbors returns true iff a and b are neighbors in the 2D matrix mat.
func AreNeighbors(mat [][lint, a, b int) bool {
	// HANDOUT: panic("TODO: implement this!")
	// BEGIN_SOLUTION
	if mat == nil {
		return false
	}
	if len(mat) == @ || len(mat[0]) == 0 {
		return false
	}
	for -, row := range mat {
		if AreAdjacent(a, b, row) {
			return true
		}
	}
	for _, row := range Transpose(mat) {
		if AreAdjacent(a, b, row) {
			return true
		}
	}
	return false
	// END_SOLUTION
}
