package sexpr

import (
	"errors"
	"math/big" // You will need to use this package in your implementation.
)

// ErrEval is the error value returned by the Evaluator if the contains
// an invalid token.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrEval = errors.New("eval error")

func (expr *SExpr) Eval() (*SExpr, error) {
	if expr.isNil() {
		// NIL is both an atom and a list.
		return mkNil(), nil
	}

	if expr.isAtom() {
		if expr.isNumber() {
			return expr, nil
		}
		if expr.isSymbol() {
			// Symbols like * and + should not be evaluated as standalone atoms.
			if expr.atom.literal == "*" || expr.atom.literal == "+" {
				return nil, ErrEval
			}
			return expr, nil
		}
	}

	if expr.isConsCell() {
		// Handle QUOTE
		if expr.car != nil && expr.car.isSymbol() && expr.car.atom.literal == "QUOTE" {
			if expr.cdr == nil || expr.cdr.isNil() {
				return nil, ErrEval
			}
			if expr.cdr.car == nil {
				return nil, ErrEval
			}
			return expr.cdr.car, nil
		}

		// Handle other functions
		if expr.car != nil && expr.car.isSymbol() {
			switch expr.car.atom.literal {
			case "CAR":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if arg.isNil() {
					// (CAR NIL) returns NIL.
					return mkNil(), nil
				}
				if !arg.isConsCell() {
					return nil, ErrEval
				}
				return arg.car, nil

			case "CDR":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if arg.isNil() {
					// (CDR NIL) returns NIL.
					return mkNil(), nil
				}
				if !arg.isConsCell() {
					return nil, ErrEval
				}
				return arg.cdr, nil

			case "CONS":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg1, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if expr.cdr.cdr == nil || expr.cdr.cdr.isNil() || expr.cdr.cdr.car == nil {
					return nil, ErrEval
				}
				arg2, err := expr.cdr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				// (CONS NIL NIL) returns (NIL . NIL).
				return mkConsCell(arg1, arg2), nil

			case "LENGTH":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if arg.isNil() {
					// (LENGTH NIL) returns 0.
					return mkNumber(big.NewInt(0)), nil
				}
				if !arg.isConsCell() {
					return nil, ErrEval
				}
				length := big.NewInt(0)
				for !arg.isNil() {
					length.Add(length, big.NewInt(1))
					if arg.cdr == nil {
						return nil, ErrEval
					}
					arg = arg.cdr
				}
				return mkNumber(length), nil

			case "+":
				if expr.cdr == nil || expr.cdr.isNil() {
					return nil, ErrEval
				}
				sum := big.NewInt(0)
				args := expr.cdr
				for !args.isNil() {
					if args.car == nil {
						return nil, ErrEval
					}
					arg, err := args.car.Eval()
					if err != nil {
						return nil, err
					}
					if !arg.isNumber() {
						return nil, ErrEval
					}
					sum.Add(sum, arg.atom.num)
					if args.cdr == nil {
						return nil, ErrEval
					}
					args = args.cdr
				}
				return mkNumber(sum), nil

			case "*":
				if expr.cdr == nil || expr.cdr.isNil() {
					return nil, ErrEval
				}
				product := big.NewInt(1)
				args := expr.cdr
				for !args.isNil() {
					if args.car == nil {
						return nil, ErrEval
					}
					arg, err := args.car.Eval()
					if err != nil {
						return nil, err
					}
					if !arg.isNumber() {
						return nil, ErrEval
					}
					product.Mul(product, arg.atom.num)
					if args.cdr == nil {
						return nil, ErrEval
					}
					args = args.cdr
				}
				return mkNumber(product), nil

			case "ATOM":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if arg.isAtom() || arg.isNil() {
					// (ATOM NIL) returns T.
					return mkSymbolTrue(), nil
				}
				return mkNil(), nil

			case "LISTP":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if arg.isConsCell() || arg.isNil() {
					// (LISTP NIL) returns T.
					return mkSymbolTrue(), nil
				}
				return mkNil(), nil

			case "ZEROP":
				if expr.cdr == nil || expr.cdr.isNil() || expr.cdr.car == nil {
					return nil, ErrEval
				}
				arg, err := expr.cdr.car.Eval()
				if err != nil {
					return nil, err
				}
				if !arg.isNumber() {
					return nil, ErrEval
				}
				if arg.atom.num.Cmp(big.NewInt(0)) == 0 {
					return mkSymbolTrue(), nil
				}
				return mkNil(), nil

			default:
				return nil, ErrEval
			}
		}
	}

	return nil, ErrEval
}