package term

import "testing"

func TestTermToString(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic in TestTermToString")
			}
		}()
		// Build a term trm
		bar := &Term{Typ: TermAtom, Literal: "bar"}
		num1 := &Term{Typ: TermNumber, Literal: "1"}
		a := &Term{Typ: TermAtom, Literal: "a"}
		foo := &Term{Typ: TermAtom, Literal: "foo"}
		X := &Term{Typ: TermVariable, Literal: "X"}
		fooX := &Term{Typ: TermCompound, Functor: foo, Args: []*Term{X}}
		trm := &Term{Typ: TermCompound, Functor: bar, Args: []*Term{num1, a, fooX}}
		// Expected string for the term trm
		expected := "bar(1, a, foo(X))"
		actual := trm.String()
		if expected != actual {
			t.Errorf("Expected %s, actual %s ", expected, actual)
		}
	}()
}
