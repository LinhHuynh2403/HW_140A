package brackets

import (
	"testing"
)

// Parser for the paired brackets language.

// Tokens for terminals and the endmarker.
func TestLexerValidTokens(t *testing.T) {
	tests := []struct {
		input       string
		expectedTok *token
	}{
		// LBRKT terminal for left bracket "(".
		{"(", mkTokenLbracket()},

		// RBRKT terminal for right bracket ")".
		{")", mkTokenRbracket()},

		// NUMBER terminal for numbers.
		{"1", mkTokenNumber("1")},
		{"-1", mkTokenNumber("-1")},
		{
			"10000000000000000000000000000000000000000000000000",
			mkTokenNumber("10000000000000000000000000000000000000000000000000"),
		},

		// The endmarker $.
		{"", mkTokenEOF()},
	}
	for idx, test := range tests {
		lex := newLexer(test.input)
		tok, err := lex.next()
		if err != nil {
			t.Errorf("\nin test %d (\"%s\"): lexer got an unexpected error %#v when tokenizing a valid input %#v", idx, test.input, err, test.input)
		}
		if !equalToken(tok, test.expectedTok) {
			t.Errorf("\nin test %d (\"%s\"):\n\texpected token %#v\n\tgot token      %#v", idx, test.input, test.expectedTok, tok)
		}
	}
}

// Strings in the bracket language:
func TestParserValid(t *testing.T) {
	for idx, input := range []string{
		"()",
		"(1)",
		"(1 2 3)",
		"(1 () ((2) 3) 4)",
	} {
		_, err := NewParser().Parse(input)
		if err != nil {
			t.Errorf("\nin test %d\nunexpected error", idx)
			continue
		}
	}
}

// Strings not in the bracket language:
func TestParserInvalid(t *testing.T) {
	for idx, test := range []string{
		// Not enclosed by a pair of top-level brackets.
		"",
		"1",
		"1 2 3 4",
		"1 (2)",
		// Contains unpaired brackets.
		"(1 (2)",
		") 1 (2",
		// Contains invalid tokens.
		"a",
		"(a)",
		"(1 a)",
	} {
		_, err := NewParser().Parse(test)
		if err == nil {
			t.Errorf("\nin test %d: should error when parsing \"%s\".", idx, test)
		}
	}
}

// LL(1) grammar for the bracket language
// <start>		 ::= <bracket>
// <bracket>	 ::= LBRKT <tail> RBRKT
// <term>		 ::= NUMBER | <bracket>
// <tail>		 ::= <term> <tail> | \epsilon

// 				|			FIRST			|			FOLLOW			||
// 				|							|							||
// <start>		|	LBRKT					|	$						||
// 				|							|							||
// <bracket>	|	LBRKT					|	LBRKT, NUMBER, RBRKT, $	||
// 				|							|							||
// <term>		|	LBRKT, NUMBER			|	LBRKT, NUMBER, RBRKT	||
// 				|							|							||
// <tail>		|	LBRKT, NUMBER, \epsilon	|	RBRKT					||
// 				|							|							||

// 				|				LBRKT				|			RBRKT			|			NUMBER			|	$	||
// 				|									|							|							|		||
// <start>		| <start> -> <bracket>				|							|							|		||
// 				|									|							|							|		||
// <bracket>	| <bracket> -> LBRKT <tail> RBRKT	|							|							|		||
// 				|									|							|							|		||
// <term>		| <term> -> <bracket>				|							| <term> -> NUMBER			|		||
// 				|									|							|							|		||
// <tail>		| <tail> -> <term> <tail>			|	<tail> -> \epsilon		| <tail> -> <term> <tail>	|		||
// 				|									|							|							|		||

// type Expr struct {
// 		Atom *token
// 		Bracket []*Expr
// }

func TestParserStructure(t *testing.T) {
	for idx, test := range []struct {
		input string;
		expected *Expr
	}{
		{"(1 (2) 3)",
		  &Expr{	// Parsed "(1 (2) 3)".
			Atom: nil,
			Bracket: []*Expr{
				&Expr{	// Parsed "1".
					Atom: mkTokenNumber("1"),
					Bracket: nil,
				},
				&Expr{	// Parsed "(2)".
					Atom: nil,
					Bracket: []*Expr{
						&Expr{	// Parsed "2".
							Atom: mkTokenNumber("2"),
							Bracket: nil,
						},
					},
				},
				&Expr{	// Parsed "3".
					Atom: mkTokenNumber("3"),
					Bracket: nil,
				},
			},
		}},
	} {
		actual, err := NewParser().Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d\nunexpected error", idx)
			continue
		}
		if actual.String() != test.expected.String() {
			t.Errorf("\nin test %d (\"%s\")\nerror: got      \"%s\"\n       expected \"%s\"",
				idx, test.input, actual.String(), test.expected.String())
		}
	}
}

func TestParser(t *testing.T) {
	for idx, test := range []struct {
		input, expectedString string
	}{
		{"( )", "()"},
		{"( 1 )", "(1)"},
		{"( 1  2 3)", "(1 2 3)"},
		{"( 1  (2) ( 3 ))", "(1 (2) (3))"},
		{"( 1 ( (2) ( 3 ( 4 )) (5)))", "(1 ((2) (3 (4)) (5)))"},
	} {
		actual, err := NewParser().Parse(test.input)
		if err != nil {
			t.Errorf("\nin test %d\nunexpected error", idx)
			continue
		}
		if actual.String() != test.expectedString {
			t.Errorf("\nin test %d (\"%s\")\nerror: got      \"%s\"\n       expected \"%s\"",
				idx, test.input, actual.String(), test.expectedString)
		}
	}
}
