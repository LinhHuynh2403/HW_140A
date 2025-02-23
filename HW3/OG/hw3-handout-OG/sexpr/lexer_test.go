package sexpr

import (
	"testing"
)

// Simple tests for tokenizing strings with a single token.
func TestLexerValidTokens(t *testing.T) {
	tests := []struct {
		input       string
		expectedTok *token
	}{
		// end-of-file, left paranthesis, right parantsis and comma tokens
		{"", mkTokenEOF()},
		{"(", mkTokenLpar()},
		{")", mkTokenRpar()},
		{".", mkTokenDot()},

		// `mkTokenQuote()` returns the special `'` notation of `QUOTE`
		{"'", mkTokenQuote()},

		// `mkTokenSymbol("QUOTE")` returns the symbol `QUOTE`
		{"QUOTE", mkTokenSymbol("QUOTE")},
		{"quoTE", mkTokenSymbol("QUOTE")},
		{"qUoTe", mkTokenSymbol("QUOTE")},
		{"quote", mkTokenSymbol("QUOTE")},

		// example of valid number tokens
		{"0", mkTokenNumber("0")},
		{"1", mkTokenNumber("1")},
		{"00001", mkTokenNumber("1")},
		{"+00001", mkTokenNumber("1")},
		{"-00001", mkTokenNumber("-1")},
		{"1234567890", mkTokenNumber("1234567890")},
		{
			"10000000000000000000000000000000000000000000000000",
			mkTokenNumber("10000000000000000000000000000000000000000000000000"),
		},
		{
			"-10000000000000000000000000000000000000000000000000",
			mkTokenNumber("-10000000000000000000000000000000000000000000000000"),
		},

		// example of valid symbol tokens
		{"+", mkTokenSymbol("+")},
		{"*", mkTokenSymbol("*")},
		{"atom", mkTokenSymbol("ATOM")},
		{"zerop", mkTokenSymbol("ZEROP")},

		{"nil", mkTokenSymbol("NIL")},
		{"nIl", mkTokenSymbol("NIL")},
		{"Nil", mkTokenSymbol("NIL")},
		{"NIL", mkTokenSymbol("NIL")},

		{"t", mkTokenSymbol("T")},
		{"T", mkTokenSymbol("T")},

		{"foo", mkTokenSymbol("FOO")},
		{"isValidAtom", mkTokenSymbol("ISVALIDATOM")},
		{"is-valid-atom", mkTokenSymbol("IS-VALID-ATOM")},
		{"random_token_0123_", mkTokenSymbol("RANDOM_TOKEN_0123_")},
		{"X", mkTokenSymbol("X")},
		{"Var1", mkTokenSymbol("VAR1")},
		{"_A_", mkTokenSymbol("_A_")},
		{"___Y__2_", mkTokenSymbol("___Y__2_")},
	}
	for idx, test := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, test.input, r)
				}
			}()
			lex := newLexer(test.input)
			tok, err := lex.next()
			if err != nil {
				t.Errorf("\nin test %d (\"%s\"): lexer got an unexpected error %#v when tokenizing a valid input %#v", idx, test.input, err, test.input)
			}
			if !equalToken(tok, test.expectedTok) {
				t.Errorf("\nin test %d (\"%s\"):\n\texpected token %#v\n\tgot token      %#v", idx, test.input, test.expectedTok, tok)
			}
		}()
	}
}

// TestLexerInvalidTokens tests that the lexer does not token invalid strings.
func TestLexerInvalidTokens(t *testing.T) {
	invalidStrings := []string{
		// Example of some invalid symbols in terms
		",",
		"\"",
		"=",
		"#",
		"$",
		"%",
		":",
	}
	for idx, input := range invalidStrings {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, input, r)
				}
			}()
			if _, err := newLexer(input).next(); err != ErrLexer {
				t.Errorf("\nin test %d (\"%s\"): lexer did not get error %#v when tokenizing an invalid input %#v",
					idx, input, ErrLexer, input)
			}
		}()
	}
}

func TestLexerSequence(t *testing.T) {
	// `newLexer(str)` returns a new lexer with given input string.
	input := " (zerop (  +  +1 '+2 (quote -3))) "
	lex := newLexer(input)
	// The expected sequence of literals when calling lex.next()
	expectedTokens := []*token{
		mkTokenLpar(),
		mkTokenSymbol("ZEROP"),
		mkTokenLpar(),
		mkTokenSymbol("+"),
		mkTokenNumber("1"),
		mkTokenQuote(),
		mkTokenNumber("2"),
		mkTokenLpar(),
		mkTokenSymbol("QUOTE"),
		mkTokenNumber("-3"),
		mkTokenRpar(),
		mkTokenRpar(),
		mkTokenRpar(),
		mkTokenEOF(),
		mkTokenEOF(),
	}
	for idx, expectedToken := range expectedTokens {
		// `lex.next()` consumes the input string, skips spaces and returns the next
		// token.
		token, err := lex.next()
		if err != nil {
			t.Errorf("lexer got an unexpected error %#v when tokenizing a valid input", err)
		}
		if token == nil {
			t.Errorf("lexer returned an unexpected nil token")
		}
		if !equalToken(token, expectedToken) {
			t.Errorf("\nin %d-th token of input \"%s\":\n\texpected token %#v\n\tgot token      %#v", idx, input, expectedToken, token)
		}
	}
}
