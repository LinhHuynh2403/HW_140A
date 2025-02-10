package sexpr

import (
	"fmt"
	"math/big"
)

// SExpr defines the struct of an S-expression.
// 1. If this S-expression is `NIL`: all `atom`, `car` and `cdr` fields should
//    be null pointers (Go's `nil`), i.e. `SExpr{}` or verbosely
//    ```
//    SExpr{
//        atom: nil,
//        car: nil,
//        cdr: nil,
//    }
//    ```
// 2. If this S-expression is a non-`NIL` atom, the `atom` field should store
//    the corresponding non-null token pointer of token types `tokenNumber` or
//    `tokenSymbol`, and both `car` and `cdr` fields should be null pointers
//    (Go's `nil`). E.g. `SExpr{atom: mkTokenSymbol("+")}` or verbosely
//    ```
//    SExpr{
//        atom: mkTokenSymbol("+"),
//        car: nil,
//        cdr: nil,
//    },
//    ```
// 3. If this S-expression is a non-`NIL` cons cell, the `atom` field should be
//    a null pointer (Go's `nil`) and both `car` and `cdr` should be non-null
//    SExpr pointers. Remember that `NIL` is represented as `&SExpr{}` but not a
//    non-pointer (Go's `nil`). E.g.
//    ```
//    SExpr{
//        atom: nil,
//        car: &SExpr{atom: mkTokenSymbol("+")},
//        cdr: &SExpr{},
//    }
//    ```
type SExpr struct {
	atom *token
	car  *SExpr
	cdr  *SExpr
}

// Below we provide useful helper functions and you may want to use them to
// create an S-expression or check whether an S-expression is a `NIL`, (symbol
// or number) atom or cons cell. __Do not modify them__. Feel free to write your
// own helper functions __in files we ask you to modify__.

// NIL
//

func mkNil() *SExpr {
	return &SExpr{}
}

// Caveat: `expr.car.isNil()` or `expr.cdr.isNil()` will not work as expected
// when `expr.isNil() == true`. Special treatments is required for this case.
func (expr *SExpr) isNil() bool {
	return expr.atom == nil && expr.car == nil && expr.cdr == nil
}

// atom
//

func mkAtom(tok *token) *SExpr {
	return &SExpr{atom: tok}
}

// Symbol, number or NIL
func (expr *SExpr) isAtom() bool {
	return expr.isNil() ||
		(expr.atom != nil && expr.car == nil && expr.cdr == nil)
}

// number atom

// Create a number atom of the int `num`
func mkNumber(num *big.Int) *SExpr {
	return &SExpr{atom: &token{typ: tokenNumber, num: num}}
}

func (expr *SExpr) isNumber() bool {
	return expr.isAtom() && !expr.isNil() && expr.atom.typ == tokenNumber
}

// symbol atom

// Create a symbol atom of the string `lit`
func mkSymbol(lit string) *SExpr {
	return &SExpr{atom: mkTokenSymbol(lit)}
}

func (expr *SExpr) isSymbol() bool {
	return expr.isAtom() && !expr.isNil() && expr.atom.typ == tokenSymbol
}

// Create a True symbol atom `T`
func mkSymbolTrue() *SExpr {
	return mkSymbol("T")
}

// cons cell
//

// Create a cons cell with given `car` and `cdr`
func mkConsCell(car, cdr *SExpr) *SExpr {
	return &SExpr{nil, car, cdr}
}

// Cons cell or NIL
func (expr *SExpr) isConsCell() bool {
	return expr.isNil() ||
		(expr.atom == nil && expr.car != nil && expr.cdr != nil)
}

// Helper functions to serialize an S-expression in different ways.

// SExprString serializes an SExpr into the __DOTTED__ S-expression representation
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
