package matrix

// If needed, you may define helper functions here.

// AreAdjacent returns true iff a and b are adjacent in lst.
<<<<<<< HEAD
func AreAdjacent(a, b int, check []int) bool {

    // The array is empty or having 1 element only
	if len(check) < 2 {
        return false
    }

    // Check if the two 
    for i, num := range check {
        if num == a {
            // Check the left of a (secure that a is not at index 0)
            if i > 0 {
                if  check[i-1] == b {
                    return true
                }
            }
            //check the right of a 
            if i < len(check)-1 {
                if check[i+1] == b {
                    return true
                }
            }
        }
    }
        // Both left and right does not have b
        return false
=======
func AreAdjacent(a, b int, lst []int) bool {
	// if the matrix length is 1 return false, since there is only one element
	if len(lst) < 2 {
		return false
	}
	/*
		first find the element a is in the matrix
		if a is in matrix, check the left and right index whether the element is b
		if the left or the right is equal to b, return true
	*/

	for i, num := range lst {
		if num == a {
			// neighbor on the right side
			if i > 0 && lst[i-1] == b {
				return true
			}
			// neighbor on the left side
			if i < len(lst)-1 && lst[i+1] == b {
				return true
			}
		}
	}
	return false
>>>>>>> main
}

// Transpose returns the transpose of the 2D matrix mat.

<<<<<<< HEAD
func Transpose(a [][]int) [][]int {
    // Check if the input matrix is empty
    if( a == nil) {return nil}
    // Check if the matrix's row or column is 0
    if len(a) == 0 || len(a[0]) == 0 {
        return [][]int{}
    }

    n := len(a)     // Number of rows in the input matrix
    m := len(a[0])  // Number of columns in the input matrix

    // Create a new 2D slice for the transposed matrix
    ans := make([][]int, m)
    for i := range ans {
        ans[i] = make([]int, n) // Each row in the transposed matrix has 'n' elements
        for j := range ans[i] {
            ans[i][j] = a[j][i] // Swap rows and columns
        }
    }

    return ans
}


func AreNeighbors(mat [][]int, a, b int) bool {
    // Check if the matrix is empty in either column or row
    if len(mat) == 0 || len(mat[0]) == 0 {
        return false
    }

    // Initialize array for use to check
    m := [4]int{-1, 1, 0, 0} 
    n := [4]int{0, 0, -1, 1}
    // Check up, down, right, and left of a
    for i := range mat {
        for j := range mat[i] {
            if mat[i][j] == a {
                for k := 0; k < 4; k++ {
                    newRow := i + m[k]
                    newCol := j + n[k]
                    if newRow >= 0 && newRow < len(mat) &&
                       newCol >= 0 && newCol < len(mat[newRow]) {
                        if mat[newRow][newCol] == b {
                            return true
                        }
                    }
                }
            }
        }
    }
    return false
}
=======
func Transpose(mat [][]int) [][]int {
	// Check if the input matrix is empty
	if mat == nil {
		return nil
	}
	if len(mat) == 0 || len(mat[0]) == 0 {
		return [][]int{}
	}

	row := len(mat)    // Number of rows in the input matrix
	col := len(mat[0]) // Number of columns in the input matrix

	// Create a new 2D slice for the transposed matrix
	tranposeMatrix := make([][]int, col)
	for i := range tranposeMatrix {
		tranposeMatrix[i] = make([]int, row) // Each row in the transpose_matrix has 'row' elements
		for j := range tranposeMatrix[i] {
			tranposeMatrix[i][j] = mat[j][i] // Swap rows and columns
		}
	}

	return tranposeMatrix
}

// find neighbor: a and b are neighbors iff they are next to each other in vertica or horizontal
func AreNeighbors(mat [][]int, a, b int) bool {
	// check the length of matrix if it is empty
	if len(mat) == 0 || len(mat[0]) == 0 {
		return false
	}
	// do the same as AreNeighbor in the first part but this time is m x n
	// find element a, and check if left, right, up, or down is element b
	// if yes, return true, otherwise, return false

	posA := [2]int{} // store location of A
	foundA := false

	// find a
	for i := range mat {
		for j := range mat[0] {
			if mat[i][j] == a {
				posA[0] = i
				posA[1] = j
				foundA = true
				break
			}
		}
		if foundA {
			break
		}
	}

	// If `a` was not found, return false
	if !foundA {
		return false
	}

	// Check the positions left, right, up, and down of `a`
	directions := [][2]int{
		{-1, 0}, // Up
		{1, 0},  // Down
		{0, -1}, // Left
		{0, 1},  // Right
	}

	for _, dir := range directions {
		newRow, newCol := posA[0]+dir[0], posA[1]+dir[1]
		if newRow >= 0 && newRow < len(mat) && newCol >= 0 && newCol < len(mat[0]) { // Ensure within bounds
			if mat[newRow][newCol] == b {
				return true
			}
		}
	}

	return false
}
>>>>>>> main
