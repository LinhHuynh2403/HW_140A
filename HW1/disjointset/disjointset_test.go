package disjointset

import (
	"math/rand"
	"testing"
)

func TestDisjointSetSimple(t *testing.T) {
	func() {
		defer func() {
			if recover() != nil {
				t.Errorf("DisjointSetSimple panicked")
			}
		}()
		s := NewDisjointSet()
		r := s.FindSet(1)
		if r != 1 {
			t.Errorf("Expected 1, actual %d", r)
		}
		r = s.UnionSet(1, 2)
		if r != 1 && r != 2 {
			t.Errorf("Expected 1 or 2, actual %d", r)
		}
		if s.FindSet(1) != s.FindSet(2) {
			t.Errorf("Expected true, actual false")
		}
	}()
}

func TestDisjointSetOddEven(t *testing.T) {
	rand.Seed(1)
	func() {
		defer func() {
			if recover() != nil {
				t.Errorf("DisjointSetOddEven panicked")
			}
		}()
		s := NewDisjointSet()
		const N = 1000 * 1000
		// Union all even numbers
		for i := 2; i < N; i += 2 {
			s.UnionSet(i, i-2)
		}
		// Union all odd numbers
		for i := 3; i < N; i += 2 {
			s.UnionSet(i, i-2)
		}
		// Perform N random checks
		for i := 1; i < N; i++ {
			j := rand.Intn(i)
			sameMod := i%2 == j%2
			sameSet := s.FindSet(i) == s.FindSet(j)
			if sameMod != sameSet {
				t.Errorf("Expected %d and %d to be in the same set", i, j)
			}
		}
	}()
}

type UnionSetStep struct {
	a, b int
}
type FindSetCheck struct {
	a, b     int
	expected bool
}

var disjointSetTests = []struct {
	unionSetSteps []UnionSetStep
	findSetChecks []FindSetCheck
}{
	{
		unionSetSteps: []UnionSetStep{},
		findSetChecks: []FindSetCheck{
			FindSetCheck{1, 1, true},
			FindSetCheck{3, 3, true},
			FindSetCheck{2, 3, false},
			FindSetCheck{6, 7, false},
		},
	},
	{
		unionSetSteps: []UnionSetStep{
			UnionSetStep{1, 2},
			UnionSetStep{1, 3},
			UnionSetStep{1, 4},
		},
		findSetChecks: []FindSetCheck{
			FindSetCheck{2, 1, true},
			FindSetCheck{2, 3, true},
			FindSetCheck{2, 5, false},
			FindSetCheck{6, 7, false},
		},
	},
	{
		unionSetSteps: []UnionSetStep{
			UnionSetStep{1, 2},
			UnionSetStep{2, 3},
			UnionSetStep{3, 4},
		},
		findSetChecks: []FindSetCheck{
			FindSetCheck{2, 4, true},
			FindSetCheck{3, 1, true},
			FindSetCheck{2, 8, false},
			FindSetCheck{3, 6, false},
		},
	},
	{
		unionSetSteps: []UnionSetStep{
			UnionSetStep{1, 2},
			UnionSetStep{2, 3},
			UnionSetStep{3, 1},
			UnionSetStep{4, 5},
			UnionSetStep{5, 6},
			UnionSetStep{6, 4},
		},
		findSetChecks: []FindSetCheck{
			FindSetCheck{1, 3, true},
			FindSetCheck{1, 6, false},
			FindSetCheck{5, 6, true},
			FindSetCheck{3, 5, false},
		},
	},
	{
		unionSetSteps: []UnionSetStep{
			UnionSetStep{1, 2},
			UnionSetStep{2, 3},
			UnionSetStep{3, 4},
			UnionSetStep{5, 6},
			UnionSetStep{6, 7},
			UnionSetStep{7, 8},
			UnionSetStep{8, 9},
		},
		findSetChecks: []FindSetCheck{
			FindSetCheck{1, 5, false},
			FindSetCheck{1, 4, true},
			FindSetCheck{5, 8, true},
		},
	},
}

func TestDisjointSet(t *testing.T) {
	testNo := 0
	for _, test := range disjointSetTests {
		testNo++
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("DisjointSet panicked on test number %d", testNo)
				}
			}()
			s := NewDisjointSet()
			for _, unionSetStep := range test.unionSetSteps {
				s.UnionSet(unionSetStep.a, unionSetStep.b)
			}
			for _, findSetCheck := range test.findSetChecks {
				if actual := s.FindSet(findSetCheck.a) == s.FindSet(findSetCheck.b); actual != findSetCheck.expected {
					t.Errorf("In test %d, FindSet(%d) == FindSet(%d) gives %t, expected %t",
						testNo, findSetCheck.a, findSetCheck.b, actual, findSetCheck.expected)
				}
			}
		}()
	}
}
