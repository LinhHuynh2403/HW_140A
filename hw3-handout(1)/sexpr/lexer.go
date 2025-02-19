package sexpr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
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
	tokenEOF tokenType = iota
	tokenSymbol
	tokenNumber
	tokenLpar
	tokenRpar
	tokenDot
	tokenQuote
)

// A token is a Lisp atom, including a number.
type token struct {
	typ     tokenType
	literal string

	// `num` is (a pointer to) an __unbounded__ integer
	// See also https://golang.org/pkg/math/big/
	num *big.Int
}

func equalToken(tok1, tok2 *token) bool {
	return tok1 != nil && tok2 != nil &&
		tok1.typ == tok2.typ &&
		tok1.literal == tok2.literal &&
		(tok1.num == tok2.num ||
			tok1.num != nil && tok2.num != nil && tok1.num.Cmp(tok2.num) == 0)
}

func (tok *token) String() string {
	if tok.typ == tokenNumber {
		return fmt.Sprintf("%d", tok.num)
	}
	return fmt.Sprintf("%s", tok.literal)
}

type lexer struct {
	rd       io.RuneReader
	peeking  bool
	peekRune rune
	last     rune
	buf      bytes.Buffer
}

func newLexer(input string) *lexer {
	return &lexer{
		rd: strings.NewReader(input),
	}
}

var tokens = make(map[string]*token)

func mkToken(typ tokenType, literal string) *token {
	tok := tokens[literal]
	if tok == nil {
		if typ == tokenNumber {
			// This error is checked in call-sites
			num, _ := new(big.Int).SetString(literal, 10)
			tok = &token{typ: typ, num: num}
		} else {
			tok = &token{typ: typ, literal: literal}
		}
	}
	return tok
}

func mkTokenEOF() *token {
	return mkToken(tokenEOF, "")
}

func mkTokenLpar() *token {
	return mkToken(tokenLpar, "(")
}

func mkTokenRpar() *token {
	return mkToken(tokenRpar, ")")
}

func mkTokenDot() *token {
	return mkToken(tokenDot, ".")
}

func mkTokenQuote() *token {
	return mkToken(tokenQuote, "QUOTE")
}

func mkTokenNumber(literal string) *token {
	return mkToken(tokenNumber, literal)
}

func mkTokenSymbol(literal string) *token {
	return mkToken(tokenSymbol, literal)
}

func (l *lexer) next() (*token, error) {
	for {
		r := l.read()
		switch {
		case isSpace(r):

		case r == eofRune:
			return mkTokenEOF(), nil

		case r == '(':
			return mkTokenLpar(), nil

		case r == ')':
			return mkTokenRpar(), nil

		case r == '.':
			return mkTokenDot(), nil

		case r == '\'':
			return mkTokenQuote(), nil

		default:
			// try to tokenize a symbol or a number
			l.accum(r, isSymbolOrNumber)
			literal := l.buf.String()

			// try to tokenize a number
			if _, ok := new(big.Int).SetString(literal, 10); ok {
				return mkTokenNumber(literal), nil
			}

			// error if neither symbol or number
			if !isSymbolOrNumber(r) {
				return nil, ErrLexer
			}

			return mkTokenSymbol(literal), nil
		}
	}
}

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
		l.buf.WriteRune(rune(unicode.ToUpper(r)))
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

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isSymbolOrNumber(r rune) bool {
	return r == '_' || r == '-' || r == '+' || r == '*' || unicode.IsLetter(r) || isNumber(r)
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}
