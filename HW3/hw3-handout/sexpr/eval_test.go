package sexpr

import (
	"testing"
)

func TestEvalInvalid(t *testing.T) {
	for idx, test := range []string{
		"x",
		"(1)",
		"(LENGTH '(1 . 2))",
		"(QUOTE)",
		"(QUOTE 1 2 3 4)",
		"(QUOTE . 1)",
		"(QUOTE . (1 . 2))",
		"(QUOTE 1 2)",
		"(QUOTE . (1 . 'NIL))",
		"(CAR)",
		"(CAR x)",
		"(CAR '(1 2) '1)",
		"(CDR)",
		"(CDR x)",
		"(CONS)",
		"(CONS 1 2 3)",
		"(CONS x 1)",
		"(CONS 1 x)",
		"(LENGTH)",
		"(LENGTH (1) ())",
		"(LENGTH 1)",
		"(LENGTH 'x)",
		"(LENGTH x)",
		"(+ x)",
		"(+ 'x)",
		"(* x)",
		"(* 'x)",
		"(ATOM)",
		"(ATOM 1 2 3)",
		"(ATOM x)",
		"(ZEROP)",
		"(ZEROP 1 2 3)",
		"(ZEROP x)",
		"(ZEROP ())",
		"(ZEROP (1))",
		"(LISTP)",
		"(LISTP 1 2 3)",
		"(LISTP x)",
		"(UNDEFINED)",
	} {
		p := NewParser()
		sexpr, err := p.Parse(test)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test)
			continue
		}
		_, err = sexpr.Eval()
		if err == nil {
			t.Errorf("\nin test %d (\"%s\"):\n\terror: should get an eval error", idx, test)
		}
	}
}

// The evaluation of symbols `*` and `+` should result in error in MiniLisp.
// Note: this behavior differs from that in CLISP.
func TestEvalArithmeticOperators(t *testing.T) {
	for idx, test := range []string{
		"*",
		"+",
	} {
		p := NewParser()
		sexpr, err := p.Parse(test)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test)
			continue
		}
		_, err = sexpr.Eval()
		if err == nil {
			t.Errorf("\nin test %d (\"%s\"):\n\terror: should get an eval error", idx, test)
		}
	}
}

func TestEvalQUOTE(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"'1", "1"},
		{"''1", "(QUOTE . (1 . NIL))"},
		{"'(1)", "(1 . NIL)"},
		{"''(1)", "(QUOTE . ((1 . NIL) . NIL))"},
		{"(QUOTE (1))", "(1 . NIL)"},
		{"(QUOTE . (1))", "1"},
		{"(QUOTE . (NIL . NIL))", "NIL"},
		{"(QUOTE . ('NIL . NIL))", "(QUOTE . (NIL . NIL))"},
		{"(QUOTE . (('1 . 2) . NIL))", "((QUOTE . (1 . NIL)) . 2)"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalNumber(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"1", "1"},
		{"+1", "1"},
		{"-001", "-1"},
		{
			"-10000000000000000000000000000000000000000000000000000000000000000",
			"-10000000000000000000000000000000000000000000000000000000000000000",
		},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalCAR(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(CAR NIL)", "NIL"},
		{"(CAR '(1 2))", "1"},
		{"(CAR '(1 . 2))", "1"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalCONS(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(CONS 1 2)", "(1 . 2)"},
		{"(CONS 1 '2)", "(1 . 2)"},
		{"(CONS 1 ''2)", "(1 . (QUOTE . (2 . NIL)))"},
		{"(CONS NIL NIL)", "(NIL . NIL)"},
		{"(CONS NIL 1)", "(NIL . 1)"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}
func TestEvalCDR(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(CDR NIL)", "NIL"},
		{"(CDR '(1 2))", "(2 . NIL)"},
		{"(CDR '(1 . 2))", "2"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}
func TestEvalLENGTH(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(LENGTH '())", "0"},
		{"(LENGTH '(1))", "1"},
		{"(LENGTH '(1 2))", "2"},
		{"(LENGTH '(1 2 3))", "3"},
		{"(LENGTH '(1 (2 3)))", "2"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalSum(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(+)", "0"},
		{"(+ 1)", "1"},
		{"(+ 1 2)", "3"},
		{"(+ 1 2 3)", "6"},
		{"(+ 1 2 3 4)", "10"},
		{"(+ 1 (+ 2 3))", "6"},
		{"(+ 1 (+ 2 3) -4)", "2"},
		{
			"(+ 1841869746456711357943187984 78943132489451238944879231278)",
			"80785002235907950302822419262",
		},
		{
			"(+ 6584453157984218784138798163 -8921871231489451877431234894)",
			"-2337418073505233093292436731",
		},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalProduct(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(*)", "1"},
		{"(* 0)", "0"},
		{"(* 1)", "1"},
		{"(* 1 2)", "2"},
		{"(* 1 2 3)", "6"},
		{"(* 1 2 3 4)", "24"},
		{"(* -1 2 3 4)", "-24"},
		{"(* -1 2 -3 4)", "24"},
		{"(* 1 (* 2 3))", "6"},
		{"(* 1 (* -2 3))", "-6"},
		{"(* 1 (* 2 3) 0)", "0"},
		{
			"(* 985617513576545648975121564 79841561894552416515616548)",
			"78693241714576626131944167058185924258870542710041072",
		},
		{
			"(* 418956489745648948124856946 -549878465467832318765413)",
			"-230375151679127070083285887312267866796671957608698",
		},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalAtom(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(ATOM NIL)", "T"},
		{"(ATOM (CAR NIL))", "T"},
		{"(ATOM (CDR NIL))", "T"},
		{"(ATOM ())", "T"},
		{"(ATOM '+1)", "T"},
		{"(ATOM '-1)", "T"},
		{"(ATOM 'some-atom)", "T"},
		{"(ATOM 1)", "T"},
		{"(ATOM '1)", "T"},
		{"(ATOM ''1)", "NIL"},
		{"(ATOM '(1))", "NIL"},
		{"(ATOM (+ 1 2 (+ 3) 4))", "T"},
		{"(ATOM '(1 . 2))", "NIL"},
		{"(ATOM (CDR '(1 . 2)))", "T"},
		{"(ATOM '(1 2))", "NIL"},
		{"(ATOM (CDR '(1 2)))", "NIL"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalLISTP(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(LISTP NIL)", "T"},
		{"(LISTP ())", "T"},
		{"(LISTP '(NIL))", "T"},
		{"(LISTP '(1))", "T"},
		{"(LISTP '(1 . 2))", "T"},
		{"(LISTP (CONS 1 2))", "T"},
		{"(LISTP (CAR NIL))", "T"},
		{"(LISTP 1)", "NIL"},
		{"(LISTP '1)", "NIL"},
		{"(LISTP ''1)", "T"},
		{"(LISTP 'x)", "NIL"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}

func TestEvalZEROP(t *testing.T) {
	for idx, test := range []struct {
		input, expected string
	}{
		{"(ZEROP 0)", "T"},
		{"(ZEROP 123456)", "NIL"},
		{"(ZEROP -999)", "NIL"},
		{"(ZEROP (+))", "T"},
		{"(ZEROP (+ (LENGTH '(1000)) -1))", "T"},
		{"(ZEROP (+ 1 (+ -1)))", "T"},
		{"(ZEROP (+ 1 2 3 4 5 (+ -15)))", "T"},
	} {
		p := NewParser()
		sexpr, err := p.Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected parse error", idx, test.input)
			continue
		}
		actual, err := sexpr.Eval()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"):\nunexpected eval error", idx, test.input)
		} else if actual.SExprString() != test.expected {
			t.Errorf("\nin test %d (\"%s\"):\nerror:\tgot\t\t\t\"%s\"\n\t\texpected\t\"%s\"",
				idx, test.input, actual.SExprString(), test.expected)
		}
	}
}
