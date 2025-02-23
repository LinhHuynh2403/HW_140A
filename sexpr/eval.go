package sexpr

import (
	"errors"
	"math/big"
)

var (
	ErrEval          = errors.New("eval error")
	ErrInvalidSyntax = errors.New("invalid syntax")
)

// Eval evaluates an S-expression.
func (expr *SExpr) Eval() (*SExpr, error) {
	if expr == nil {
		return nil, ErrEval
	}

	// Handle atoms (numbers, symbols, and NIL)
	if expr.isAtom() {
		if expr.isSymbol() && expr.atom.literal != "NIL" && expr.atom.literal != "T" {
			return nil, ErrEval // Symbols other than NIL and T are invalid
		}
		return expr, nil
	}

	// Handle cons cells (lists and special forms)
	if !expr.isConsCell() {
		return nil, ErrEval
	}

	// Evaluate the CAR of the list (the function or special form)
	fn := expr.car
	if !fn.isSymbol() {
		return nil, ErrEval
	}

	// Evaluate special forms and functions
	switch fn.atom.literal {
	case "QUOTE":
		return evalQuote(expr.cdr)
	case "CAR":
		return evalCar(expr.cdr)
	case "CDR":
		return evalCdr(expr.cdr)
	case "CONS":
		return evalCons(expr.cdr)
	case "LENGTH":
		return evalLength(expr.cdr)
	case "+":
		return evalSum(expr.cdr)
	case "*":
		return evalProduct(expr.cdr)
	case "ATOM":
		return evalAtom(expr.cdr)
	case "LISTP":
		return evalListp(expr.cdr)
	case "ZEROP":
		return evalZerop(expr.cdr)
	default:
		return nil, ErrEval // Unknown function or special form
	}
}

// evalQuote handles the QUOTE special form.
func evalQuote(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	return args.car, nil
}

// evalCar handles the CAR function.
func evalCar(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	list, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	if !list.isConsCell() {
		return nil, ErrEval
	}
	return list.car, nil
}

// evalCdr handles the CDR function.
func evalCdr(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	list, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	if !list.isConsCell() {
		return nil, ErrEval
	}
	return list.cdr, nil
}

// evalCons handles the CONS function.
func evalCons(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isConsCell() || !args.cdr.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	car, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	cdr, err := args.cdr.car.Eval()
	if err != nil {
		return nil, err
	}
	return mkConsCell(car, cdr), nil
}

// evalLength handles the LENGTH function.
func evalLength(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	list, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	length := big.NewInt(0)
	for !list.isNil() {
		if !list.isConsCell() {
			return nil, ErrEval
		}
		length.Add(length, big.NewInt(1))
		list = list.cdr
	}
	return mkNumber(length), nil
}

// evalSum handles the + function.
func evalSum(args *SExpr) (*SExpr, error) {
	sum := big.NewInt(0)
	for args != nil && args.isConsCell() {
		arg, err := args.car.Eval()
		if err != nil {
			return nil, err
		}
		if !arg.isNumber() {
			return nil, ErrEval
		}
		sum.Add(sum, arg.atom.num)
		args = args.cdr
	}
	return mkNumber(sum), nil
}

// evalProduct handles the * function.
func evalProduct(args *SExpr) (*SExpr, error) {
	product := big.NewInt(1)
	for args != nil && args.isConsCell() {
		arg, err := args.car.Eval()
		if err != nil {
			return nil, err
		}
		if !arg.isNumber() {
			return nil, ErrEval
		}
		product.Mul(product, arg.atom.num)
		args = args.cdr
	}
	return mkNumber(product), nil
}

// evalAtom handles the ATOM predicate.
func evalAtom(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	arg, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	if arg.isAtom() {
		return mkSymbolTrue(), nil
	}
	return mkNil(), nil
}

// evalListp handles the LISTP predicate.
func evalListp(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	arg, err := args.car.Eval()
	if err != nil {
		return nil, err
	}
	if arg.isConsCell() || arg.isNil() {
		return mkSymbolTrue(), nil
	}
	return mkNil(), nil
}

// evalZerop handles the ZEROP predicate.
func evalZerop(args *SExpr) (*SExpr, error) {
	if args == nil || !args.isConsCell() || !args.cdr.isNil() {
		return nil, ErrInvalidSyntax
	}
	arg, err := args.car.Eval()
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
