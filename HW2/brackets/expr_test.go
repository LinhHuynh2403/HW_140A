package brackets

import (
	"math/big"
	"testing"
)

func TestExprNumber(t *testing.T) {
	expr := mkAtom(mkTokenNumber("100"))
	if !expr.isAtom() {
		t.Errorf("expected isAtom() to return true")
	}
	if expr.isBracket() {
		t.Errorf("expected isBracket() to return false")
	}
	if expr.String() != "100" {
		t.Errorf("incorrect String()")
	}
}

func TestExprBracket(t *testing.T) {
	expr := mkBracket([]*Expr{
		mkNumber(big.NewInt(1)),
		mkBracket([]*Expr{
			mkNumber(big.NewInt(2)),
		}),
		mkNumber(big.NewInt(3)),
	})
	if expr.isAtom() {
		t.Errorf("expected isAtom() to return false")
	}
	if !expr.isBracket() {
		t.Errorf("expected isBracket() to return true")
	}
	if expr.String() != "(1 (2) 3)" {
		t.Errorf("incorrect String()")
	}
}
