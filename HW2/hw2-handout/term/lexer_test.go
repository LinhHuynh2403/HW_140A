package term

import "testing"

// Simple tests for tokenizing strings with a single token.
func TestLexerValidTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected tokenType
	}{
		// end-of-file, left paranthesis, right parantsis and comma tokens
		{"", tokenEOF},
		{"(", tokenLpar},
		{")", tokenRpar},
		{",", tokenComma},
		// example of valid number tokens
		{"0", tokenNumber},
		{"1", tokenNumber},
		{"1234567890", tokenNumber},
		// example of valid atom tokens
		{"f", tokenAtom},
		{"foo", tokenAtom},
		{"isValidAtom", tokenAtom},
		{"is_valid_atom", tokenAtom},
		{"random_token_0123_", tokenAtom},
		// example of valid variable tokens
		{"X", tokenVariable},
		{"Var1", tokenVariable},
		{"_A_", tokenVariable},
		{"___Y__2_", tokenVariable},
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
			if tok.literal != test.input {
				t.Errorf("\nin test %d (\"%s\"): expected token literal %#v, got token literal %#v", idx, test.input, test.input, tok.literal)
			}
		}()
	}
}

// TestLexerInvalidTokens tests that the lexer does not token invalid strings.
func TestLexerInvalidTokens(t *testing.T) {
	invalidStrings := []string{
		// Example of some invalid symbols in terms
		"'",
		"\"",
		"+",
		"-",
		"*",
		"=",
		"#",
		"$",
		"%",
		":",
		// Invalid atom tokens: we do not consider quoted atoms
		"'X'",
		// Invalid number tokens
		"0123",
		"123a",
		"1a23",
		"+123",
		"-123",
		// Invalid variable token: we do not consider the wildcard variable "_"
		"_",
	}
	for idx, input := range invalidStrings {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("\nin test %d (\"%s\") panic: %s", idx, input, r)
				}
			}()
			if _, err := newLexer(input).next(); err != ErrLexer {
				t.Errorf("\nin test %d (\"%s\"): lexer did not get error %#v when tokenizing an invalid input %#v", idx, input, err, input)
			}
		}()
	}
}

func TestLexerSequence(t *testing.T) {
	// `newLexer(str)` returns a new lexer with given input string.
	lex := newLexer(" foo ( 1 ) ")
	// The expected sequence of literals when calling lex.next()
	expectedTokens := []Token{
		{tokenAtom, "foo"},
		{tokenLpar, "("},
		{tokenNumber, "1"},
		{tokenRpar, ")"},
		{tokenEOF, ""},
		{tokenEOF, ""},
	}
	for _, expectedToken := range expectedTokens {
		// `lex.next()` consumes the input string, skips spaces and returns the next
		// token.
		token, err := lex.next()
		if err != nil {
			t.Errorf("lexer got an unexpected error %#v when tokenizing a valid input", err)
		}
		if token == nil {
			t.Errorf("lexer returned an unexpected nil token")
		}
		if token.typ != expectedToken.typ {
			t.Errorf("expected token type %#v, got token type %#v", expectedToken.typ, token.typ)
		}
		if token.literal != expectedToken.literal {
			t.Errorf("expected token literal %#v, got token literal %#v", expectedToken.literal, token.literal)
		}
	}
}
