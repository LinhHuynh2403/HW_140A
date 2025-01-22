package matrix

// If needed, you may define helper functions here.

// AreAdjacent returns true iff a and b are adjacent in lst.
func AreAdjacent(a, b int, check []int) bool {
	panic("TODO: implement this!")
	if len(check) < 2 {
        return false
    }
//1-1-2-1-1
    for i, num := range check {
        if num == a {
            if i > 0 && check[i-1] == b {
                return true
            }
            if i < len(check)-1 && check[i+1] == b {
                return true
            }
        }
    }
        return false
}

// Transpose returns the transpose of the 2D matrix mat.
func Transpose(a [][]int) [][]int {
	panic("TODO: implement this!")
	if len(a) == 0 {
        return [][]int{} 
    }
    m := len(a[0])   //colum
    n := len(a)     //row
    ans := make([][]int, m)
    for i := range ans {
        ans[i] = make([]int, n)
        for j := range ans[i] {
            ans[i][j] = a[j][i]
        }
    }
    return ans
}


func AreNeighbors(mat [][]int, a, b int) bool {
    if len(mat) == 0 || len(mat[0]) == 0 {
        return false
    }
    m := [4]int{-1, 1, 0, 0} 
    n := [4]int{0, 0, -1, 1}
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
