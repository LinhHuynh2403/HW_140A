package sexpr

import (
	"math/big"
	"testing"
)

func buildNilExample() *SExpr {
	return mkNil()
}

func buildNumberExample() *SExpr {
	return mkNumber(big.NewInt(100))
}

func buildSymbolExample() *SExpr {
	return mkSymbol("+")
}

func buildConsCellExample() *SExpr {
	return mkConsCell(mkSymbol("A"), mkSymbol("B"))
}

func buildProperListExample() *SExpr {
	return mkConsCell(mkSymbol("A"), mkConsCell(mkSymbol("B"), mkConsCell(mkSymbol("C"), mkNil())))
}

func buildDottedListExample() *SExpr {
	return mkConsCell(mkSymbol("A"), mkConsCell(mkSymbol("B"), mkSymbol("C")))
}

func buildQuoteExample() *SExpr {
	return mkConsCell(mkSymbol("QUOTE"), mkConsCell(mkSymbol("A"), mkNil()))
}

type SExprBuilder func() *SExpr

func TestParseExample(t *testing.T) {
	for idx, test := range []struct {
		input                string
		expectedSExprBuilder SExprBuilder
	}{
		{"()", buildNilExample},
		{"100", buildNumberExample},
		{"+", buildSymbolExample},
		{"(A . B)", buildConsCellExample},
		{"(A B C)", buildProperListExample},
		{"(A B . C)", buildDottedListExample},
		{"'A", buildQuoteExample},
	} {
		actual, err := NewParser().Parse(test.input)
		expected := test.expectedSExprBuilder()
		if err != nil {
			t.Errorf("\nin test %d\nunexpected error", idx)
			continue
		}
		if actual.SExprString() != expected.SExprString() {
			t.Errorf("\nin test %d (\"%s\")\nerror: got      %s\n       expected %s",
				idx, test.input, actual.SExprString(), expected.SExprString())
		}
	}
}

func TestParserInvalid(t *testing.T) {
	for idx, test := range []string{
		"",
		"(",
		"'",
		")",
		"x)",
		"( ) ( )",
		"(a . () . () . ())",
		"((x .",
		"(x",
	} {
		_, err := NewParser().Parse(test)
		if err == nil {
			t.Errorf("\nin test %d\nshould error", idx)
		}
	}
}

func TestParserProperList(t *testing.T) {
	for idx, test := range []struct {
		input, expectedSExprString string
	}{
		{"()", "NIL"},
		{"a", "A"},

		// proper lists
		{"(())", "(NIL . NIL)"},
		{"(a)", "(A . NIL)"},
		{"((a))", "((A . NIL) . NIL)"},
		{"(a b c)", "(A . (B . (C . NIL)))"},
		{"(a b c d)", "(A . (B . (C . (D . NIL))))"},
		{"(  a b c  d  e)", "(A . (B . (C . (D . (E . NIL)))))"},
		{"(  (a b )c)", "((A . (B . NIL)) . (C . NIL))"},
		{"(a b (c d))", "(A . (B . ((C . (D . NIL)) . NIL)))"},
		{"(a () () a)", "(A . (NIL . (NIL . (A . NIL))))"},

		// dotted lists
		{"(a (b . c))", "(A . ((B . C) . NIL))"},
		{"(a . b)", "(A . B)"},
		{"(a . (b . c))", "(A . (B . C))"},
		{"(a . (b . (c . d)))", "(A . (B . (C . D)))"},
		{"(a . ((b . c) . d))", "(A . ((B . C) . D))"},
		{"(a b . c)", "(A . (B . C))"},
		{"(a b c . d)", "(A . (B . (C . D)))"},
		{"(a (b c) d . e)", "(A . ((B . (C . NIL)) . (D . E)))"},
		{"(a b (c d) . e)", "(A . (B . ((C . (D . NIL)) . E)))"},
		{"(a b c . (d e))", "(A . (B . (C . (D . (E . NIL)))))"},
		{"(a b c . (d . e))", "(A . (B . (C . (D . E))))"},
		{"(a(b.c  ))", "(A . ((B . C) . NIL))"},
		{"(a . ( (  ) . ( ( ) . a)))", "(A . (NIL . (NIL . A)))"},
	} {
		actual, err := NewParser().Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d\nunexpected error", idx)
			continue
		}
		if actual.SExprString() != test.expectedSExprString {
			t.Errorf("\nin test %d (\"%s\")\nerror: got      %s\n       expected %s",
				idx, test.input, actual.SExprString(), test.expectedSExprString)
		}
	}
}

func TestParseQuote(t *testing.T) {
	for idx, test := range []struct {
		input, expectedSExpr string
	}{
		{"'(1 2)", "(QUOTE . ((1 . (2 . NIL)) . NIL))"},
		{"'(1 . 2)", "(QUOTE . ((1 . 2) . NIL))"},
		{"(quote . (1 . 2))", "(QUOTE . (1 . 2))"},
		{"'a", "(QUOTE . (A . NIL))"},
		{"'(a)", "(QUOTE . ((A . NIL) . NIL))"},
		{"''a", "(QUOTE . ((QUOTE . (A . NIL)) . NIL))"},
		{"''(a)", "(QUOTE . ((QUOTE . ((A . NIL) . NIL)) . NIL))"},
		{"(' a 'b '  c)", "((QUOTE . (A . NIL)) . ((QUOTE . (B . NIL)) . ((QUOTE . (C . NIL)) . NIL)))"},
	} {
		actual, _ := NewParser().Parse(test.input)
		if actual.SExprString() != test.expectedSExpr {
			t.Errorf("\nerror: in test %d (\"%s\"):\n\texpected: %s\n\tgot      %s",
				idx, test.input, test.expectedSExpr, actual.SExprString())
		}
	}
}
