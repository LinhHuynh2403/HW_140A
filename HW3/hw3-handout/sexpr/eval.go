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
			return evalNil(expr)
	}
	if expr.isAtom() {
			return evalAtom(expr)
	}
	if expr.isConsCell() {
			return evalConsCell(expr)
	}
	return nil, ErrEval
}

func evalNil(expr *SExpr) (*SExpr, error) {
	if expr.isNil() {
		// NIL is both an atom and a list.
		return mkNil(), nil
	}
	return nil, ErrEval
}
	
func evalAtom(expr *SExpr) (*SExpr, error) {
	if expr.isNumber() {
		return expr, nil
	}
	if expr.isSymbol() {
		// Symbols like * and + should not be evaluated as standalone atoms.
		if expr.atom.literal == "*" || expr.atom.literal == "+" {
			return nil, ErrEval
		}
	}
	
	return nil, ErrEval
}

func evalConsCell(expr *SExpr) (*SExpr, error) {
	if expr.car != nil && expr.car.isSymbol() {
		switch expr.car.atom.literal {
		case "QUOTE":
			return evalQUOTE(expr)
		case "CAR":
			return evalCAR(expr)
		case "CDR":
			return evalCDR(expr)
		case "CONS":
			return evalCONS(expr)
		case "LENGTH":
			return evalLENGTH(expr)
		case "ATOM":
			return evalATOM(expr)
		case "LISTP":
			return evalLISTP(expr)
		case "ZEROP":
			return evalZEROP(expr)
		case "+":
			return evalSUM(expr)
		case "*":
			return evalPRODUCT(expr)
		}
	}
	
	return nil, ErrEval
}

func evalQUOTE(expr *SExpr) (*SExpr, error) {
	if expr.cdr == nil || expr.cdr.isConsCell() {
		return nil, ErrEval
	}
	if expr.cdr.car == nil {
		return nil, ErrEval
	}
	return expr.cdr.car, nil
}
func evalCAR(expr *SExpr) (*SExpr, error) {
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
}

func evalCDR(expr *SExpr) (*SExpr, error) {
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
}

func evalCONS(expr *SExpr) (*SExpr, error) {
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
}
func evalLENGTH(expr *SExpr) (*SExpr, error) {
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
}

func evalATOM(expr *SExpr) (*SExpr, error) {
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
}
func evalLISTP(expr *SExpr) (*SExpr, error) {
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
}

func evalZEROP(expr *SExpr) (*SExpr, error) {
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
}

func evalSUM(expr *SExpr) (*SExpr, error) {
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
}

func evalPRODUCT(expr *SExpr) (*SExpr, error) {
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
}