package term

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

// ErrLexer is the error value returned by the Lexer if the contains
// an invalid token.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrLexer = errors.New("lexer error")

// tokenType enumerates all types to tokens
// See also https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go
type tokenType int

const (
	tokenEOF      tokenType = iota // for the End-of-File (EoF) token
	tokenLpar                      // for the left paranthesis token '('
	tokenRpar                      // for the right parenthesis token ')'
	tokenComma                     // for the comma token ','
	tokenAtom                      // for the atom tokens like 'foo123', which start with a lower case letter.
	tokenNumber                    // for the (natural) number tokens like '0' and '123'
	tokenVariable                  // for the variable tokens like 'Var2' and '_X', which start with a capital letter or underscore.
)

// Token struct
type Token struct {
	typ     tokenType // contains the type of this token
	literal string    // contains the literal of this token, e.g. '123', 'foo', ')'
}

// Lexer struct
type lexer struct {
	rd       io.RuneReader
	last     rune
	peeking  bool
	peekRune rune
	buf      bytes.Buffer
	// tokens stores the existing token pointers for the corresponding token
	// literals
	tokens map[string]*Token
}

// newLexer creates a new instance of the type lexer with the input string
func newLexer(input string) *lexer {
	return &lexer{
		rd:     strings.NewReader(input),
		tokens: make(map[string]*Token),
	}
}

// The rune represents EoF
const eofRune rune = -1

func (l *lexer) read() rune {
	if l.peeking {
		l.peeking = false
		return l.peekRune
	}
	r, _, err := l.rd.ReadRune()
	if err == io.EOF {
		r = eofRune
	}
	l.last = r
	return r
}

func (l *lexer) accum(r rune, valid func(rune) bool) {
	l.buf.Reset()
	for {
		l.buf.WriteRune(r)
		r = l.read()
		if r == eofRune {
			return
		}
		if !valid(r) {
			l.back(r)
			return
		}
	}
}

func (l *lexer) back(r rune) {
	l.peeking = true
	l.peekRune = r
}

func (l *lexer) peek() rune {
	r := l.read()
	l.back(r)
	return r
}

// mkToken returns the unique token of the given token type and
// literal
func (l *lexer) mkToken(typ tokenType, literal string) *Token {
	tok, ok := l.tokens[literal]
	if !ok {
		tok = &Token{typ, literal}
		l.tokens[literal] = tok
	}
	return tok
}

// l.next() consumes the input string of this lexer, skips spaces and and return
// (nextToken, nil) if the next token is valid, otherwise returns (nil, ErrLexer).
func (l *lexer) next() (*Token, error) {
	for {
		// read the next rune
		r := l.read()
		switch {
		// skip if the next rune is space
		case isSpace(r):
		case r == eofRune:
			return l.mkToken(tokenEOF, ""), nil
		case r == '(':
			return l.mkToken(tokenLpar, "("), nil
		case r == ')':
			return l.mkToken(tokenRpar, ")"), nil
		case r == ',':
			return l.mkToken(tokenComma, ","), nil

		case isNumberZero(r):
			if !l.nextRuneIsSeparator() {
				return nil, ErrLexer
			}
			return l.mkToken(tokenNumber, "0"), nil

		case isNumberNonZero(r):
			l.accum(r, isNumber)
			if !l.nextRuneIsSeparator() {
				return nil, ErrLexer
			}
			return l.mkToken(tokenNumber, l.buf.String()), nil

		case isAtomHead(r):
			l.accum(r, isAlphaNumUnderscore)
			return l.mkToken(tokenAtom, l.buf.String()), nil

		case isVariableHead(r):
			l.accum(r, isAlphaNumUnderscore)
			if l.buf.String() == "_" {
				return nil, ErrLexer // We do not consider the wildcard variable "_"
			}
			return l.mkToken(tokenVariable, l.buf.String()), nil

		default:
			return nil, ErrLexer
		}
	}
}

// Check if the rune r can separate the last token from the next token (if exists)
func (l *lexer) nextRuneIsSeparator() bool {
	r := l.peek()
	return isSeparator(r)
}

func isSeparator(r rune) bool {
	return isSpace(r) || r == eofRune || r == '(' || r == ')' || r == ','
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isNumberZero(r rune) bool {
	return r == '0'
}

func isNumberNonZero(r rune) bool {
	return '0' <= r && r <= '9'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAtomHead(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func isVariableHead(r rune) bool {
	return r == '_' || ('A' <= r && r <= 'Z')
}

func isAlphaNumUnderscore(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || isNumber(r)
}
