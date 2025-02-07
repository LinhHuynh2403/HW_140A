package term

import (
	"errors"
	// "strconv"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <start>    ::= <term> | \epsilon
// <term>     ::= ATOM | NUM | VAR | <compound>
// <compound> ::= <functor> LPAR <args> RPAR
// <functor>  ::= ATOM
// <args>     ::= <term> | <term> COMMA <args>
//

// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

// termParser implements the Parser interface.
type termParser struct {
	lexer *lexer
	memo  map[string]*Term
}

// NewParser creates a new instance of the termParser.
func NewParser() Parser {
	return &termParser{memo: make(map[string]*Term)}
}

// Parse parses the input string and returns a DAG representation of the term.
func (p *termParser) Parse(input string) (*Term, error) {
	// Handle empty input (ε case)
	if input == "" {
		return nil, nil
	}

	p.lexer = newLexer(input)
	tok, err := p.lexer.next()
	if err != nil {
		return nil, ErrParser
	}

	return p.parseTerm(tok)
}

// parseTerm parses a term according to the grammar.
func (p *termParser) parseTerm(tok *Token) (*Term, error) {
	switch tok.typ {
	case tokenAtom:
		return p.getOrCreateTerm(TermAtom, tok.literal, nil), nil
	case tokenNumber:
		return p.getOrCreateTerm(TermNumber, tok.literal, nil), nil
	case tokenVariable:
		return p.getOrCreateTerm(TermVariable, tok.literal, nil), nil
	case tokenLpar:
		return p.parseCompound()
	default:
		return nil, ErrParser
	}
}

// parseCompound parses a compound term.
func (p *termParser) parseCompound() (*Term, error) {
	// Expect a functor (ATOM)
	tok, err := p.lexer.next()
	if err != nil || tok.typ != tokenAtom {
		return nil, ErrParser
	}
	functor := tok.literal

	// Expect '('
	tok, err = p.lexer.next()
	if err != nil || tok.typ != tokenLpar {
		return nil, ErrParser
	}

	// Parse arguments
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	// Expect ')'
	tok, err = p.lexer.next()
	if err != nil || tok.typ != tokenRpar {
		return nil, ErrParser
	}

	// Return compound term, ensuring DAG structure
	return p.getOrCreateTerm(TermCompound, functor, args), nil
}

// parseArgs parses the arguments of a compound term.
func (p *termParser) parseArgs() ([]*Term, error) {
	var args []*Term
	for {
		tok, err := p.lexer.next()
		if err != nil {
			return nil, ErrParser
		}
		term, err := p.parseTerm(tok)
		if err != nil {
			return nil, err
		}
		args = append(args, term)

		tok, err = p.lexer.next()
		if err != nil {
			return nil, ErrParser
		}

		// If we hit ')', we are done parsing arguments
		if tok.typ == tokenRpar {
			p.lexer.back(rune(tok.literal[0])) // Push back to reprocess ')'
			break
		}

		// Expect ','
		if tok.typ != tokenComma {
			return nil, ErrParser
		}
	}
	return args, nil
}

// getOrCreateTerm ensures unique DAG representation by memoizing sub-terms.
func (p *termParser) getOrCreateTerm(typ TermType, literal string, args []*Term) *Term {
	// Create a unique key for the term
	key := literal
	if typ == TermCompound {
		for _, arg := range args {
			key += ":" + arg.Literal
		}
	}

	// If already exists in memo, return existing reference
	if existing, found := p.memo[key]; found {
		return existing
	}

	// Otherwise, create a new term and store it in memo
	newTerm := &Term{Typ: typ, Literal: literal, Args: args}
	p.memo[key] = newTerm
	return newTerm
}
