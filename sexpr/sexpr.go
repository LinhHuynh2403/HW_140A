package sexpr

import (
	"fmt"
	"math/big"
)

// SExpr defines the structure of an S-expression.
type SExpr struct {
	atom *token
	car  *SExpr
	cdr  *SExpr
}

// Creates an empty S-expression representing NIL
func mkNil() *SExpr {
	return &SExpr{}
}

// Checks if the S-expression is NIL
func (expr *SExpr) isNil() bool {
	return expr.atom == nil && expr.car == nil && expr.cdr == nil
}

// Creates an atom S-expression from a token
func mkAtom(tok *token) *SExpr {
	return &SExpr{atom: tok}
}

// Checks if the S-expression is an atom (symbol, number, or NIL)
func (expr *SExpr) isAtom() bool {
	return expr.isNil() || (expr.atom != nil && expr.car == nil && expr.cdr == nil)
}

// Creates a number atom S-expression
func mkNumber(num *big.Int) *SExpr {
	return &SExpr{atom: &token{typ: tokenNumber, num: num}}
}

// Checks if the S-expression is a number atom
func (expr *SExpr) isNumber() bool {
	return expr.isAtom() && !expr.isNil() && expr.atom.typ == tokenNumber
}

// Creates a symbol atom S-expression
func mkSymbol(lit string) *SExpr {
	return &SExpr{atom: mkTokenSymbol(lit)}
}

// Checks if the S-expression is a symbol atom
func (expr *SExpr) isSymbol() bool {
	return expr.isAtom() && !expr.isNil() && expr.atom.typ == tokenSymbol
}

// Creates a True symbol atom "T"
func mkSymbolTrue() *SExpr {
	return mkSymbol("T")
}

// Creates a cons cell S-expression
func mkConsCell(car, cdr *SExpr) *SExpr {
	return &SExpr{nil, car, cdr}
}

// Checks if the S-expression is a cons cell (or NIL)
func (expr *SExpr) isConsCell() bool {
	return expr.isNil() || (expr.atom == nil && expr.car != nil && expr.cdr != nil)
}

// Serializes an S-expression into the dotted representation
func (expr *SExpr) SExprString() string {
	switch {
	case expr.isNil():
		return "NIL"
	case expr.isAtom():
		return expr.atom.String()
	default:
		return fmt.Sprintf("(%s . %s)", expr.car.SExprString(), expr.cdr.SExprString())
	}
}
