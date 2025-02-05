package branch

import (
	// "fmt"
	"testing"
)

//!+ComputeBranchFactors
func TestComputeBranchFactors(t *testing.T) {
	var test_code = `
	package main

	import (
		"fmt"
		"eval"
	)

	func sequential_if() {
		if true {
			fmt.Println("is true")
		} else {
			fmt.Println("isn't true")
		}

		if true {
			fmt.Println("is true")
		} else {
			fmt.Println("isn't true")
		}

		if true {
			fmt.Println("is true")
		} else {
			fmt.Println("isn't true")
		}

		if true {
			fmt.Println("is true")
		} else {
			fmt.Println("isn't true")
		}
	}

	func single_for() {
		for i := 0; i < 10; i += 1 {
			fmt.Println("is loop")
		}

		for i := 0; i < 10; i += 1 {
			fmt.Println("is loop")
		}

		for i := 0; i < 10; i += 1 {
			fmt.Println("is loop")
		}
	}

	func single_range() {
		for i := range []int{1, 1, 2, 3, 5, 8} {
			fmt.Println("is fibonacci")
		}

		for i := range []int{1, 1, 2, 3, 5, 8} {
			fmt.Println("is fibonacci")
		}

		for i := range []int{1, 1, 2, 3, 5, 8} {
			fmt.Println("is fibonacci")
		}

		for i := range []int{1, 1, 2, 3, 5, 8} {
			fmt.Println("is fibonacci")
		}

		for i := range []int{1, 1, 2, 3, 5, 8} {
			fmt.Println("is fibonacci")
		}
	}

	func no_branches() {
		return 42
	}

	func single_switch() {
		switch 5 {
		case 0:
			// pass
		case 5:
			fmt.Println("It's five!")
		default:
			fmt.Println("It isn't five...")
		}

		switch 5 {
		case 0:
			// pass
		case 5:
			fmt.Println("It's five!")
		default:
			fmt.Println("It isn't five...")
		}
	}

	func single_typeswitch() {
		var x interface{} = "test"
		switch x.(type) {
		case uint:
			// pass
		case string:
			fmt.Println("It's a string!")
		default:
			fmt.Println("It's not a string...")
		}
		switch x.(type) {
		case uint:
			// pass
		case string:
			fmt.Println("It's a string!")
		default:
			fmt.Println("It's not a string...")
		}
		switch x.(type) {
		case uint:
			// pass
		case string:
			fmt.Println("It's a string!")
		default:
			fmt.Println("It's not a string...")
		}
	}

	func nested_if(x int) float64 {
		var result float64
		if x < 0 {
			if x > -5 {
				result := -0.5
			} else {
				result := -1
			}
		} else if x > 0 {
			result := 1
		} else {
			result := 0
		}
		return result
	}

	func nested_for_if() {
		for i := 0; i < 10; i += 1 {
			if i > 5 {
				fmt.Println("is filter")
			}
		}
	}

	func nested_switch_if(x int) {
		switch x > 5 {
		case true:
			if x > 10 {
				fmt.Println("is really big")
			}
		default:
			fmt.Println("is default")
		}
	}

	func mixed_switch_no_default_for_if() {
		switch 5 {
		case 0:
			// pass
		case 5:
			fmt.Println("It's five!")
		}

		for i := 0; i < 10; i += 1 {
			if i > 5 {
				fmt.Println("is filter")
			}
		}
	}

	func single_typeswitch_no_default() {
		var x interface{} = "test"
		switch x.(type) {
		case uint:
			// pass
		case string:
			fmt.Println("It's a string!")
		}
	}

	func nested_if_no_else(x int) float64 {
		var result float64 = 0
		if x < 0 {
			if x > -5 {
				result := -0.5
			} else {
				result := -1
			}
		} else if x > 0 {
			result := 1
		}
		return result
	}
	`

	tests := []struct {
		name     string
		branches uint
	}{
		{"sequential_if", 4},
		{"single_for", 3},
		{"single_range", 5},
		{"no_branches", 0},
		{"single_switch", 2},
		{"single_typeswitch", 3},
		{"nested_if", 3},
		{"nested_for_if", 2},
		{"nested_switch_if", 2},
		{"mixed_switch_no_default_for_if", 3},
		{"single_typeswitch_no_default", 1},
		{"nested_if_no_else", 3},
	}

	branch_factors := ComputeBranchFactors(test_code)

	for _, test := range tests {
		if branch_factors[test.name] != test.branches {
			t.Errorf("ComputeBranchFactors(%v) = %d, want %d\n",
				test.name, branch_factors[test.name], test.branches)
		}
	}
}

//!-ComputeBranchFactors

func TestComputeBranchFactors_Fail(t *testing.T) {

	func() {

		defer func() {
			if recover() == nil {
				t.Errorf("did not panic, but should\n")
			}
		}()
		ComputeBranchFactors("not a valid go program")
	}()

}
