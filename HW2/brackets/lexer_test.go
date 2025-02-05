package brackets

import (
	"testing"
)

// Simple tests for tokenizing strings with a single token.
func TestLexerValid(t *testing.T) {
	tests := []struct {
		input       string
		expectedTok *token
	}{
		// end-of-file, left paranthesis, right parantsis and comma tokens
		{"", mkTokenEOF()},
		{"(", mkTokenLbracket()},
		{")", mkTokenRbracket()},

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
			if tok.String() != test.expectedTok.String() || !equalToken(tok, test.expectedTok) {
				t.Errorf("\nin test %d (\"%s\"):\n\texpected token %#v\n\tgot token      %#v", idx, test.input, test.expectedTok, tok)
			}
		}()
	}
}

// TestLexerInvalidTokens tests that the lexer does not token invalid strings.
func TestLexerInvalid(t *testing.T) {
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
	input := " ( (  1 2 (3))) "
	lex := newLexer(input)
	// The expected sequence of literals when calling lex.next()
	expectedTokens := []*token{
		mkTokenLbracket(),
		mkTokenLbracket(),
		mkTokenNumber("1"),
		mkTokenNumber("2"),
		mkTokenLbracket(),
		mkTokenNumber("3"),
		mkTokenRbracket(),
		mkTokenRbracket(),
		mkTokenRbracket(),
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
