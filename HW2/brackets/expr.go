package brackets

import (
	"fmt"
	"math/big"
)

// Expr defines a struct of an expression of the brackets langauge.
//
// 1. If the expression is an atom expression, the `Atom` field contains the
// corresponding token pointer of token types `tokenNumber`, and the `Bracket`
// field should be a null pointer.
//
// 2. If the expression is an bracket expression, the `Atom` field should be a
// null pointer and the `Bracket` field should point to a slice of the
// expressions inside the bracket.
//
// For example, the bracket expression `(1 (2 3))` is represented as
// ```
// 	Expr {										// Expr of `(1 (2 3))`.
//		Atom: nil,
//		Bracket: []*Expr{
//			&Expr{								// Expr of `1`.
//				Atom: mkTokenNumber(1),
//				Brackets: nil
//			},
//			&Expr{								// Expr of `(2 3)`.
//				Atom: nil,
//				Brackets: []*Expr{
//					&Expr{						// Expr of `2`.
//						Atom: mkTokenNumber(2),
//						Brackets: nil
//					},
//					&Expr{						// Expr of `3`.
//						Atom: mkTokenNumber(3),
//						Brackets: nil
//					},
//				}
//			},
//		}
// 	}
// ```

type Expr struct {
	Atom *token
	Bracket []*Expr
}

// atom
//

// Make an atom experssion from a token.
func mkAtom(tok *token) *Expr {
	return &Expr{Atom: tok}
}

// Check if the expression is an atom expression.
func (expr *Expr) isAtom() bool {
	return expr.Atom != nil && expr.Bracket == nil
}

// Create a number atom of the int `num`
func mkNumber(num *big.Int) *Expr {
	return &Expr{Atom: &token{typ: tokenNumber, num: num}}
}

// Make a bracket experssion from a slice of inner expressions.
func mkBracket(bracket []*Expr) *Expr {
	return &Expr{Bracket: bracket}
}

// Check if an experssion is a bracket expression
func (expr *Expr) isBracket() bool {
	return expr.Atom == nil && expr.Bracket != nil
}

// Serializes a bracket expersion.
func (expr *Expr) String() string {
	switch {
	case expr.isAtom():
		return expr.Atom.String()
	default:
		str := "("
		for _, e := range expr.Bracket {
			if str == "(" {
				str = fmt.Sprintf("%s%s", str, e.String())
			} else {
				str = fmt.Sprintf("%s %s", str, e.String())
			}
		}
		str = fmt.Sprintf("%s)", str)
		return str
	}
}
