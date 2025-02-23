package sexpr

import (
	"math/big"
	"testing"
)

func TestSExprNil(t *testing.T) {
	se := mkNil()
	if !se.isNil() {
		t.Errorf("expected isNil() to return true")
	}
	if !se.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if !se.isConsCell() {
		t.Errorf("expected isConsCell() to return true")
	}
	if se.isNumber() {
		t.Errorf("expected isNumber() to return false")
	}
	if se.isSymbol() {
		t.Errorf("expected isSymbol() to return false")
	}
	if se.SExprString() != "NIL" {
		t.Errorf("incorrect SExprString()")
	}
}

func TestSExprAtom(t *testing.T) {
	se := mkAtom(mkTokenSymbol("A"))
	if se.isNil() {
		t.Errorf("expected isNil() to return false")
	}
	if !se.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if se.isConsCell() {
		t.Errorf("expected isConsCell() to return false")
	}
	if se.SExprString() != "A" {
		t.Errorf("incorrect SExprString()")
	}
}

func TestSExprNumber(t *testing.T) {
	se := mkNumber(big.NewInt(100))
	if se.isNil() {
		t.Errorf("expected isNil() to return false")
	}
	if !se.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if se.isConsCell() {
		t.Errorf("expected isConsCell() to return false")
	}
	if !se.isNumber() {
		t.Errorf("expected isNumber() to return true")
	}
	if se.isSymbol() {
		t.Errorf("expected isSymbol() to return false")
	}
	if se.SExprString() != "100" {
		t.Errorf("incorrect SExprString()")
	}
}

func TestSExprSymbol(t *testing.T) {
	se := mkSymbol("+")
	if se.isNil() {
		t.Errorf("expected isNil() to return false")
	}
	if !se.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if se.isConsCell() {
		t.Errorf("expected isConsCell() to return false")
	}
	if se.isNumber() {
		t.Errorf("expected isNumber() to return false")
	}
	if !se.isSymbol() {
		t.Errorf("expected isSymbol() to return true")
	}
	if se.SExprString() != "+" {
		t.Errorf("incorrect SExprString()")
	}
}

func TestSExprSymbolTrue(t *testing.T) {
	se := mkSymbolTrue()
	if se.isNil() {
		t.Errorf("expected isNil() to return false")
	}
	if !se.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if se.isConsCell() {
		t.Errorf("expected isConsCell() to return false")
	}
	if se.isNumber() {
		t.Errorf("expected isNumber() to return false")
	}
	if !se.isSymbol() {
		t.Errorf("expected isSymbol() to return true")
	}
	if se.SExprString() != "T" {
		t.Errorf("incorrect SExprString()")
	}
}

func TestSExprConsCell(t *testing.T) {
	se := mkConsCell(mkSymbol("A"), mkSymbol("B"))
	if se.isNil() {
		t.Errorf("expected isNil() to return false")
	}
	if se.isAtom() {
		t.Errorf("expected isAtom() to return false")
	}
	if !se.isConsCell() {
		t.Errorf("expected isConsCell() to return true")
	}
	if se.isNumber() {
		t.Errorf("expected isNumber() to return false")
	}
	if se.isSymbol() {
		t.Errorf("expected isSymbol() to return false")
	}
	if se.SExprString() != "(A . B)" {
		t.Errorf("incorrect SExprString()")
	}
}
