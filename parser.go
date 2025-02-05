package term

import (
	"errors"
	"fmt"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
var ErrParser = errors.New("parser error")

// Token types for the lexer
type tokenType int

const (
	tokenEOF tokenType = iota
	tokenLpar
	tokenRpar
	tokenComma
	tokenAtom
	tokenNumber
	tokenVariable
)

// Token struct for lexer tokens
type Token struct {
	typ     tokenType
	literal string
}

// TermType enumerates all types of terms
type TermType int

const (
	TermAtom TermType = iota
	TermNumber
	TermVariable
	TermCompound
)

// Term represents a term in the grammar
type Term struct {
	Typ     TermType
	Literal string
	Functor *Term
	Args    []*Term
}

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
var ErrParser = errors.New("parser error")

// Parser is the interface for the term parser.
type Parser interface {
	Parse(string) (*Term, error)
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	return &parser{}
}

type parser struct {
	lexer *lexer
	token *Token
}

func (p *parser) Parse(input string) (*Term, error) {
	p.lexer = newLexer(input)
	p.token, _ = p.lexer.next()

	// Parse the start symbol <start> ::= <term> | ε
	term, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	return term, nil
}

// Parse the <term> ::= ATOM | NUM | VAR | <compound>
func (p *parser) parseTerm() (*Term, error) {
	switch p.token.typ {
	case tokenAtom:
		return p.parseAtom()
	case tokenNumber:
		return p.parseNumber()
	case tokenVariable:
		return p.parseVariable()
	case tokenLpar:
		return p.parseCompound()
	default:
		return nil, ErrParser
	}
}

// Parse the ATOM type term
func (p *parser) parseAtom() (*Term, error) {
	atom := &Term{
		Typ:     TermAtom,
		Literal: p.token.literal,
	}
	p.token, _ = p.lexer.next()
	return atom, nil
}

// Parse the NUM type term
func (p *parser) parseNumber() (*Term, error) {
	number := &Term{
		Typ:     TermNumber,
		Literal: p.token.literal,
	}
	p.token, _ = p.lexer.next()
	return number, nil
}

// Parse the VAR type term
func (p *parser) parseVariable() (*Term, error) {
	variable := &Term{
		Typ:     TermVariable,
		Literal: p.token.literal,
	}
	p.token, _ = p.lexer.next()
	return variable, nil
}

// Parse a <compound> term
func (p *parser) parseCompound() (*Term, error) {
	// <compound> ::= <functor> LPAR <args> RPAR
	functor, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	// Expecting left parenthesis after the functor
	if p.token.typ != tokenLpar {
		return nil, ErrParser
	}
	p.token, _ = p.lexer.next()

	// Parse arguments (could be a single term or multiple terms)
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	// Expecting right parenthesis after arguments
	if p.token.typ != tokenRpar {
		return nil, ErrParser
	}
	p.token, _ = p.lexer.next()

	// Return the compound term
	return &Term{
		Typ:     TermCompound,
		Functor: functor,
		Args:    args,
	}, nil
}

// Parse arguments in the compound term
func (p *parser) parseArgs() ([]*Term, error) {
	var args []*Term
	for {
		arg, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		// Look ahead for comma or right parenthesis
		if p.token.typ == tokenComma {
			p.token, _ = p.lexer.next()
		} else {
			break
		}
	}
	return args, nil
}
