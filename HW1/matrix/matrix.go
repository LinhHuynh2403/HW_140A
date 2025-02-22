package matrix

// If needed, you may define helper functions here.

// AreAdjacent returns true iff a and b are adjacent in lst.
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
}

// Transpose returns the transpose of the 2D matrix mat.

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