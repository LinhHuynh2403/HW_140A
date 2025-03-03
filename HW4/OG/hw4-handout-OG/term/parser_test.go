package term

import (
	"fmt"
	"testing"
)

func TestParserInvalidTerms(t *testing.T) {
	for idx, input := range []string{
		// Invalid terms
		"f)",
		"f()",
		"f((",
		"f(1)g",
		",f(1)",
		"f(1),",
		"f(X",
		"(X, 1)",
		"X, 1)",
		", 1)",
		"F(X)",
		"123(X)",
		// TODO add more tests for 100% test coverage
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, input, r)
				}
			}()
			p := NewParser()
			if _, err := p.Parse(input); err == nil {
				t.Errorf("parser did not got error %#v when parsing an invalid input %#v", err, input)
			}
		}()
	}
}

type termTestGeneratorFunction func() (string, *Term)

func termWithoutSharingTest0() (string, *Term) {
	return "", nil
}

func termWithoutSharingTest1() (string, *Term) {
	return "0", &Term{Typ: TermNumber, Literal: "0"}
}

func termWithoutSharingTest2() (string, *Term) {
	return "123", &Term{Typ: TermNumber, Literal: "123"}
}

func termWithoutSharingTest3() (string, *Term) {
	return "foo", &Term{Typ: TermAtom, Literal: "foo"}
}

func termWithoutSharingTest4() (string, *Term) {
	return "sizeOf", &Term{Typ: TermAtom, Literal: "sizeOf"}
}

func termWithoutSharingTest5() (string, *Term) {
	return "X", &Term{Typ: TermVariable, Literal: "X"}
}

func termWithoutSharingTest6() (string, *Term) {
	return "_X_1", &Term{Typ: TermVariable, Literal: "_X_1"}
}

func termWithoutSharingTest7() (string, *Term) {
	f := &Term{Typ: TermAtom, Literal: "f"}
	X := &Term{Typ: TermVariable, Literal: "X"}
	return "f(X)", &Term{Typ: TermCompound, Functor: f, Args: []*Term{X}}
}

func termWithoutSharingTest8() (string, *Term) {
	foo := &Term{Typ: TermAtom, Literal: "foo"}
	a := &Term{Typ: TermAtom, Literal: "a"}
	X := &Term{Typ: TermVariable, Literal: "X"}
	return "foo  ( a , X )", &Term{Typ: TermCompound, Functor: foo, Args: []*Term{a, X}}
}

func termWithoutSharingTest9() (string, *Term) {
	bar := &Term{Typ: TermAtom, Literal: "bar"}
	num1 := &Term{Typ: TermNumber, Literal: "1"}
	a := &Term{Typ: TermAtom, Literal: "a"}
	foo := &Term{Typ: TermAtom, Literal: "foo"}
	X := &Term{Typ: TermVariable, Literal: "X"}
	fooX := &Term{Typ: TermCompound, Functor: foo, Args: []*Term{X}}
	return "bar(1 ,a, foo( X ))", &Term{Typ: TermCompound, Functor: bar, Args: []*Term{num1, a, fooX}}
}

func termWithoutSharingTest10() (string, *Term) {
	return "f(A,g(B,h(C,D)),p(E))", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "f"},
		Args: []*Term{
			&Term{Typ: TermVariable, Literal: "A"},
			&Term{
				Typ:     TermCompound,
				Functor: &Term{Typ: TermAtom, Literal: "g"},
				Args: []*Term{
					&Term{Typ: TermVariable, Literal: "B"},
					&Term{
						Typ:     TermCompound,
						Functor: &Term{Typ: TermAtom, Literal: "h"},
						Args: []*Term{
							&Term{Typ: TermVariable, Literal: "C"},
							&Term{Typ: TermVariable, Literal: "D"},
						},
					},
				},
			},
			&Term{
				Typ:     TermCompound,
				Functor: &Term{Typ: TermAtom, Literal: "p"},
				Args: []*Term{
					&Term{Typ: TermVariable, Literal: "E"},
				},
			},
		},
	}
}

func TestParseTermWithoutSharing(t *testing.T) {
	for idx, testGenerator := range []termTestGeneratorFunction{
		termWithoutSharingTest0,
		termWithoutSharingTest1,
		termWithoutSharingTest2,
		termWithoutSharingTest3,
		termWithoutSharingTest4,
		termWithoutSharingTest5,
		termWithoutSharingTest6,
		termWithoutSharingTest7,
		termWithoutSharingTest8,
		termWithoutSharingTest9,
		termWithoutSharingTest10,
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					input, _ := testGenerator()
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, input, r)
				}
			}()
			input, expected := testGenerator()
			p := NewParser()
			actual, err := p.Parse(input)
			if err != nil {
				t.Errorf("\nin test %d (\"%s\") parser got unexpected error: %#v", idx, input, err)
			}
			if areIsomorphic, err := checkIsomorphic(expected, actual); !areIsomorphic {
				t.Errorf("\nin test %d (\"%s\")%s", idx, input, err)
			}
		}()
	}
}

func termWithSharingTest0() (string, *Term) {
	X := &Term{Typ: TermVariable, Literal: "X"}
	return "rel(X, X)", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "rel"},
		Args: []*Term{
			X,
			X,
		},
	}
}

func termWithSharingTest1() (string, *Term) {
	X := &Term{Typ: TermVariable, Literal: "X"}
	return "foo  ( X ,X, X)  ", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "foo"},
		Args: []*Term{
			X,
			X,
			X,
		},
	}
}

func termWithSharingTest2() (string, *Term) {
	X := &Term{Typ: TermVariable, Literal: "X"}
	return " foo( X, X ,f (X) )", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "foo"},
		Args: []*Term{
			X,
			X,
			&Term{
				Typ:     TermCompound,
				Functor: &Term{Typ: TermAtom, Literal: "f"},
				Args: []*Term{
					X,
				}},
		},
	}
}

func termWithSharingTest3() (string, *Term) {
	f := &Term{Typ: TermAtom, Literal: "f"}
	X := &Term{Typ: TermVariable, Literal: "X"}
	fX := &Term{Typ: TermCompound, Functor: f, Args: []*Term{X}}
	return "foo ( X, X , X, f(X), f(f (X) ))", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "foo"},
		Args: []*Term{
			X,
			X,
			X,
			fX,
			&Term{
				Typ:     TermCompound,
				Functor: f,
				Args: []*Term{
					fX,
				}},
		},
	}
}

func termWithSharingTest4() (string, *Term) {
	fX := &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "f"},
		Args: []*Term{
			&Term{Typ: TermVariable, Literal: "X"},
		},
	}
	return "rel( f( X ) , f (X) )", &Term{
		Typ:     TermCompound,
		Functor: &Term{Typ: TermAtom, Literal: "rel"},
		Args: []*Term{
			fX,
			fX,
		},
	}
}

func TestParseTermWithSharing(t *testing.T) {
	for idx, testGenerator := range []termTestGeneratorFunction{
		termWithSharingTest0,
		termWithSharingTest1,
		termWithSharingTest2,
		termWithSharingTest3,
		termWithSharingTest4,
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					input, _ := testGenerator()
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, input, r)
				}
			}()
			p := NewParser()
			input, expected := testGenerator()
			actual, err := p.Parse(input)
			if err != nil {
				t.Errorf("\nin test %d (\"%s\") parser got unexpected error: %#v", idx, input, err)
			}
			if areIsomorphic, err := checkIsomorphic(expected, actual); !areIsomorphic {
				t.Errorf("\nin test %d (\"%s\")%s", idx, input, err)
			}
		}()
	}
}

func checkIsomorphic(expected, actual *Term) (bool, error) {
	matchTerms := make(map[*Term]*Term)
	return checkTermIsomorphic(expected, actual, matchTerms)
}

func checkTermIsomorphic(expected, actual *Term, matchTerms map[*Term]*Term) (bool, error) {
	if x, ok := matchTerms[expected]; ok {
		if x == actual {
			return true, nil
		}
		return false, fmt.Errorf(
			"\nerror:\n|\tthe subterm:\n|\t\t\"%s\" (%#v(%p))\n|\tin the expected term matches more than one terms:\n|\t\t\"%s\" (%#v(%p))\n|\t\t\"%s\" (%#v(%p))\n|\tin the actual term",
			expected, expected, expected,
			x, x, x,
			actual, actual, actual)
	}
	if expected != actual {
		if (expected == nil || actual == nil) ||
			(expected.Typ != actual.Typ) ||
			(expected.Literal != actual.Literal) {
			return false, fmt.Errorf(
				"\nerror:\n|\texpected\n|\t\t\"%s\" (%#v(%p))\n|\tgot\n|\t\t\"%s\" (%#v(%p))",
				expected, expected, expected,
				actual, actual, actual)
		}
		if areIsomorphic, err := checkTermIsomorphic(expected.Functor, actual.Functor, matchTerms); !areIsomorphic {
			return false, fmt.Errorf(
				"\ncontext:\n|\tin the functor of\n|\t\t\"%s\" (%#v(%p))\n|\tand\n|\t\t\"%s\" (%#v(%p))%s",
				expected, expected, expected,
				actual, actual, actual,
				err)
		}
		if areIsomorphic, err := checkTermSliceIsomorphic(expected.Args, actual.Args, matchTerms); !areIsomorphic {
			return false, fmt.Errorf(
				"\ncontext:\n|\tin the arguments of\n|\t\t\"%s\" (%#v(%p))\n|\tand\n|\t\t\"%s\" (%#v(%p))%s",
				expected, expected, expected,
				actual, actual, actual,
				err)
		}
	}
	matchTerms[expected] = actual
	return true, nil
}

func checkTermSliceIsomorphic(expectedSlice, actualSlice []*Term, matchTerms map[*Term]*Term) (bool, error) {
	if (expectedSlice == nil && actualSlice != nil) ||
		(expectedSlice != nil && actualSlice == nil) ||
		(len(expectedSlice) != len(actualSlice)) {
		return false, fmt.Errorf(
			"\nerror:\n|\texpected:\n|\t\t\"(%s)\" (%#v(%p))\n|\tgot:\n|\t\t\"(%s)\" (%#v(%p))",
			TermSliceToString(expectedSlice), expectedSlice, expectedSlice,
			TermSliceToString(actualSlice), actualSlice, actualSlice)
	}
	for idx := range expectedSlice {
		if areIsomorphic, err := checkTermIsomorphic(expectedSlice[idx], actualSlice[idx], matchTerms); !areIsomorphic {
			return false, fmt.Errorf("\n|\tin the %d-th argument:%s", idx+1, err)
		}
	}
	return true, nil
}
