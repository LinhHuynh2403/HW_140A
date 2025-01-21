package matrix

// If needed, you may define helper functions here.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AreAdjacent returns true iff a and b are adjacent in lst.
func AreAdjacent(a, b int, lst []int) bool {
	panic("TODO: implement this!")
	// to check if a and b are adjacent, we check the index of a and b
	// if (abs(index(a) - index(b) == 1)) then return true, otherwise, return false
	var index_a, index_b = -1, -1

	for i := 0; i < len(lst); i++ {
		if lst[i] == a {
			index_a = i
		}
		if lst[i] == b {
			index_b = i
		}
	}
	// when a and b are not in the matrix
	if index_a == -1 || index_b == -1 {
		return false
	}
	return abs(index_a-index_b) == 1
}

// Transpose returns the transpose of the 2D matrix mat.
func Transpose(mat [][]int) [][]int {
	panic("TODO: implement this!")
	// create a new transpose_matrix n x m matrix
	var m, n = len(mat), len(mat[0])

	// create new result matrix (slice)
	transpose_matrix := make([][]int, n)
	for i := range transpose_matrix {
		transpose_matrix[i] = make([]int, m)
	}

	// loop through the origin array
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			transpose_matrix[j][i] = mat[i][j]
		}
	}

	return transpose_matrix
}

// AreNeighbors returns true iff a and b are neighbors in the 2D matrix mat.
func AreNeighbors(mat [][]int, a, b int) bool {
	panic("TODO: implement this!")
	// a and b are neighbors iff they are next to each other in vertical or horizontal
	// iterate through matrix and find location of a and b
	// a(i1, j1) and b(i2, j2)
	var i1, j1, i2, j2 int
	var m, n = len(mat), len(mat[0])

	i1, i2, j1, j2 = -1, -1, -1, -1 // when a and b are not in the matrix

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] == a {
				i1 = i
				j1 = j
			}
			if mat[i][j] == b {
				i2 = i
				j2 = j
			}
		}
	}

	if ((abs(i1-i2) == 1) && j1 == j2) || ((abs(j1-j2) == 1) && i1 == i2) {
		return true
	}
	return false

}
