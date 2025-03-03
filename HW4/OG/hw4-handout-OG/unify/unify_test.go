package unify

import (
	"hw4/term"
	"reflect"
	"testing"
)

// Unification failure due to symbol clash
func TestUnifyErrorSymbolClash(t *testing.T) {
	testCases := []struct {
		input1, input2 string
	}{
		{"f", "g"},
		{"1", "2"},
		{"f", "1"},
		{"f(X)", "g(X)"},
		{"f(X, Y)", "f(1)"},
		{"f(1, 2)", "f(A)"},
		{
			"foo(A, B, C, D, E, F, G, H)",
			"foo(A, B, C, D, E, F, G, H, I)",
		},
		{
			"f1(1, X, f2(Y, foo), f3(bar), Y, 99)",
			"f1(1, X, f2(Y, foo), f3(bar, Y), 99)",
		},
		{"f(X, s)", "f(X, 1)"},
		{"f(g(1), g(X))", "f(g(s), g(2))"},
		{"f(g(1), g(X))", "h"},
		{"f(X, g(Y), h(g(Z)))", "f(Y, f(Z), X)"},
		{
			"f1(f2(Y, f3(X, 1, f4(Z, f5(99, s1, s2), s1), s2), X))",
			"f1(f2(2, f3(s3, 1, f4(X, f5(99, s1, s2), Y), s2), X))",
		},
	}
	for idx, test := range testCases {
		unifier := NewUnifier()
		parser := term.NewParser()
		term1, err := parser.Parse(test.input1)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input1)
			continue
		}
		term2, err := parser.Parse(test.input2)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input2)
			continue
		}
		_, err = unifier.Unify(term1, term2)
		if err == nil {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\tdid not get error", idx, test.input1, test.input2)
		}
	}
}

// Unification failure due to a cycle (occurs check fails)
func TestUnifyErrorExistsACycle(t *testing.T) {
	testCases := []struct {
		input1, input2 string
	}{
		{"X", "f(X)"},
		{"f(X, f(Y))", "f(Y, X)"},
		{"a(b(c, 2, X), 0, Z)", "X"},
		{
			"f(A, B, C, D, E, F)",
			"f(f(B), C, D, E, F, A)",
		},
		{
			"f(f(A), f(B), f(C), f(D), f(E), f(F))",
			"f(B, C, D, E, F, A)",
		},
		{
			"f(f(F), D, f(D), f(A), F, f(B))",
			"f(f(E), f(C), E, B, A, C)",
		},
		{
			"f(f(A), f(f(B)), f(f(C)), f(f(f(D))), f(f(f(f(E)))), f(f(f(f(f(F))))))",
			"f(B, C, D, E, F, A)",
		},
		{
			"f(g(A), h(B, C), f(g(D, 1)), h(A, g(E)))",
			"f(E, h(D, A), f(g(C, 1)), h(A, g(B)))",
		},
	}
	for idx, test := range testCases {
		unifier := NewUnifier()
		parser := term.NewParser()
		term1, err := parser.Parse(test.input1)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input1)
			continue
		}
		term2, err := parser.Parse(test.input2)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input2)
			continue
		}
		_, err = unifier.Unify(term1, term2)
		if err == nil {
			t.Errorf("in test %d when unifying %#v and %#v: did not get error", idx, test.input1, test.input2)
		}
	}
}

func TestUnifySuccess(t *testing.T) {
	for idx, test := range []struct {
		input1, input2 string
	}{
		{"f(X, X)", "f(A, B)"},
		{"f(X, Y)", "f(g(Y), 2)"},
		{"f(A, B, C)", "f(X, X, X)"},
		{"f(h(Y), Y)", "f(X, g(1))"},
		{"f(X, g(a))", "f(g(Y), g(Y))"},
		{"f(X, Y)", "f(f(Y), f(Z))"},
		{"f(X, Y, Z)", "f(2, f(Z), f(X))"},
		{
			"f(A, B, C, D, E, F, 1)",
			"f(X, X, X, X, X, X, X)",
		},
		{
			"f(A, B, C, D, E, F, 1)",
			"f(X, X, X, X, X, X, C)",
		},
		{
			"f(A, B, C, D, E, F, 1)",
			"f(X, A, B, C, D, E, F)",
		},
		{
			"f(A, B, f(C), f(g(D, E)), f(g(h(F), 2)), g(1))",
			"f(X, A, f(B), f(g(C, D)), f(g(h(E), 2)), g(F))",
		},
		{
			"f(A, B, C, D, E, 2)",
			"f(C, A, 1, F, D, E)",
		},
	} {
		unifier := NewUnifier()
		parser := term.NewParser()
		term1, err := parser.Parse(test.input1)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input1)
			continue
		}
		term2, err := parser.Parse(test.input2)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input2)
			continue
		}
		_, err = unifier.Unify(term1, term2)
		if err != nil {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\tgot unexpected error", idx, test.input1, test.input2)
			continue
		}
	}
}

type UnifyResultAsStr map[string]string

// ToStringMap converts `UnifyResult` to the string form
func (result UnifyResult) ToStringMap() UnifyResultAsStr {
	stringMap := UnifyResultAsStr{}
	for k, v := range result {
		stringMap[k.String()] = v.String()
	}
	return stringMap
}

func TestUnifyUnique(t *testing.T) {
	for idx, test := range []struct {
		input1, input2   string
		expectedAsStrMap UnifyResultAsStr
	}{
		{"X", "1", UnifyResultAsStr{"X": "1"}},
		{"Y", "f", UnifyResultAsStr{"Y": "f"}},
		{"X", "f(1)", UnifyResultAsStr{"X": "f(1)"}},
		{"f(Y)", "X", UnifyResultAsStr{"X": "f(Y)"}},
		{"f(X)", "f(1)", UnifyResultAsStr{"X": "1"}},
		{"f(X, 1)", "f(2, Y)", UnifyResultAsStr{"X": "2", "Y": "1"}},
		{"f(X, g(1))", "f(2, Y)", UnifyResultAsStr{"X": "2", "Y": "g(1)"}},
		{"f(h(Z), Y)", "f(X, g(1))", UnifyResultAsStr{"X": "h(Z)", "Y": "g(1)"}},
		{"a(b(c, 2, X), 0, Z)", "A", UnifyResultAsStr{"A": "a(b(c, 2, X), 0, Z)"}},
		{
			"f(A, B, C, D, E, F)",
			"f(1, 2, 3, 4, 5, 6)",
			UnifyResultAsStr{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5", "F": "6"},
		},
	} {
		unifier := NewUnifier()
		parser := term.NewParser()
		term1, err := parser.Parse(test.input1)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input1)
			continue
		}
		term2, err := parser.Parse(test.input2)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input2)
			continue
		}
		actual, err := unifier.Unify(term1, term2)
		if err != nil {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\tgot unexpected error", idx, test.input1, test.input2)
			continue
		}
		if actualAsStrMap := actual.ToStringMap(); !reflect.DeepEqual(test.expectedAsStrMap, actualAsStrMap) {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\texpected: %#v\n\tgot     : %#v",
				idx, test.input1, test.input2, test.expectedAsStrMap, actualAsStrMap)
		}
	}
}

// Test case where the MGU is not unique but can be a small number of possibilities
func TestUnifyNotUnique(t *testing.T) {
	for idx, test := range []struct {
		input1, input2  string
		possibleResults []UnifyResultAsStr
	}{
		{
			"X",
			"Y",
			[]UnifyResultAsStr{
				{"X": "Y"},
				{"Y": "X"},
			},
		},
		{
			"f(X, g(a))",
			"f(g(Y), g(Y))",
			[]UnifyResultAsStr{
				{"X": "g(Y)", "Y": "a"},
				{"X": "g(a)", "Y": "a"},
			},
		},
	} {
		unifier := NewUnifier()
		parser := term.NewParser()
		term1, err := parser.Parse(test.input1)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input1)
			continue
		}
		term2, err := parser.Parse(test.input2)
		if err != nil {
			t.Errorf("\nin test %d when parsing %#v:\n\tgot unexpected error", idx, test.input2)
			continue
		}
		result, err := unifier.Unify(term1, term2)
		if err != nil {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\tgot unexpected error", idx, test.input1, test.input2)
			continue
		}

		resultAsStr := result.ToStringMap()
		resultIsCorrect := false
		for _, possibleResult := range test.possibleResults {
			if reflect.DeepEqual(resultAsStr, possibleResult) {
				resultIsCorrect = true
				break
			}
		}

		if !resultIsCorrect {
			t.Errorf("\nin test %d when unifying %#v and %#v:\n\tpossible results are: %#v\n\tgot: %#v",
				idx, test.input1, test.input2, test.possibleResults, resultAsStr)
		}
	}
}
