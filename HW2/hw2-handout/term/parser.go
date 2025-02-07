package term

import (
	"errors"
 	"strconv"
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

// LL(1) grammar for the given expressions
// <start>      ::= <term> | ϵ
// <term>       ::= ATOM <pars> | NUM | VAR
// <pars>       ::= (<args>) | ϵ
// <args>       ::= <term> <other args>
// <other args> ::= , <args> | ϵ

//              |            FIRST            |          FOLLOW          ||
//              |                             |                          ||
// <start>      | ATOM, NUM, VAR, ϵ           | $                        ||
//              |                             |                          ||
// <term>       | ATOM, NUM, VAR              | ), $                     ||
//              |                             |                          ||
// <pars>       | (, ϵ                        | ), $                     ||
//              |                             |                          ||
// <args>       | ATOM, NUM, VAR              | )                        ||
//              |                             |                          ||
// <other args> | ',', ϵ                      | ), $                     ||
//              |                             |                          ||

//              |      ATOM                      |      NUM       				   |      VAR       				|      (        	  |      )        		|      ,        			|      $        	||
//              |                                |                				   |                				|               	  |               		|               			|               	||
// <start>      | <start> -> <term>              | <start> -> <term>  			   | <start> -> <term>  			| <start> -> ϵ  	  |               		|               			| <start> -> ϵ  	||
//              |                                |                				   |                				|               	  |               		|               			|               	||
// <term>       | <term> -> ATOM <pars>          | <term> -> NUM  				   | <term> -> VAR  				|               	  |               		|               			|               	||
//              |                                |                				   |                				|               	  |               		|               			|               	||
// <pars>       |                                |                				   |                				| <pars> -> (<args>)  | <pars> -> ϵ			|               			| <pars> -> ϵ   	||
//              |                                |    							   |                				|               	  |               		|               			|               	||
// <args>       | <args> -> <term> <other args>  | <args> -> <term> <other args>   | <args> -> <term> <other args>  |                	  |               		|               			|               	||
//              |                                |                				   |                				|               	  |               		|               			|               	||
// <other args> |                                |                				   |                				|               	  | <other args> -> ϵ   | <other args> -> , <args>  | <other args> -> ϵ ||
//              |               				 |                				   |                				|               	  |               		|               			|               	||


// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

// Implement the Parser interface.
type ParserImpl struct {
	lex     *lexer
	peekTok *Token
	usedTerms map[string]*Term
}


// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	return &ParserImpl{usedTerms: make(map[string]*Term)}
}

// Helper function which returns the next Token.
func (p *ParserImpl) nextToken() (*Token, error) {
	if tok := p.peekTok; tok != nil {
		p.peekTok = nil
		return tok, nil
	}

	tok, err := p.lex.next()
	if err != nil {
		return nil, ErrParser
	}

	return tok, nil
}

// Helper function which puts a Token back as the next Token.
func (p *ParserImpl) backToken(tok *Token) {
	p.peekTok = tok
}

// Helper function to peek the next Token.
func (p *ParserImpl) peekToken() (*Token, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}

	p.backToken(tok)

	return tok, nil
}

// Implement the Parse method of the Parser interface for ParserImpl.
// Parse parses the input string and returns a Term or an error.
func (p *ParserImpl) Parse(input string) (*Term, error) {
	p.lex = newLexer(input)
	p.peekTok = nil

	// Start parsing from the <start> non-terminal
	term, err := p.parseStart()
	if err != nil {
		return nil, err
	}

	// Ensure the input is fully consumed
	tok, err := p.nextToken()
	if err != nil || tok.typ != tokenEOF {
		return nil, ErrParser
	}

	return term, nil
}

// parseStart parses the <start> non-terminal.
func (p *ParserImpl) parseStart() (*Term, error) {
	tok, err := p.peekToken()
	if err != nil {
		return nil, err
	}

	// <start> ::= <term> | ϵ
	if tok.typ == tokenEOF {
		return nil, nil // Epsilon production
	}

	return p.parseTerm()
}

// parseTerm parses the <term> non-terminal.
func (p *ParserImpl) parseTerm() (*Term, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	// <term> ::= ATOM <pars> | NUM | VAR
	switch tok.typ {
	case tokenAtom:
		// ATOM <pars>
		compound, err := p.parsePars()
		if err != nil {
			return nil, err
		}
		if compound != nil {
			functor := p.getOrCreateTerm(TermAtom, tok.literal)
			return p.getOrCreateCompoundTerm(functor, compound.Args), nil
		}
		return p.getOrCreateTerm(TermAtom, tok.literal), nil
	case tokenNumber:
		// NUM
		return p.getOrCreateTerm(TermNumber, tok.literal), nil
	case tokenVariable:
		// VAR
		return p.getOrCreateTerm(TermVariable, tok.literal), nil
	default:
		return nil, ErrParser
	}
}

// parsePars parses the <pars> non-terminal.
func (p *ParserImpl) parsePars() (*Term, error) {
	tok, err := p.peekToken()
	if err != nil {
		return nil, err
	}

	// <pars> ::= (<args>) | ϵ
	if tok.typ != tokenLpar {
		return nil, nil // Epsilon production
	}

	// Consume '('
	_, err = p.nextToken()
	if err != nil {
		return nil, err
	}

	// Parse <args>
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	// Expect ')'
	tok, err = p.nextToken()
	if err != nil || tok.typ != tokenRpar {
		return nil, ErrParser
	}

	return &Term{Typ: TermCompound, Args: args}, nil
}

// parseArgs parses the <args> non-terminal.
func (p *ParserImpl) parseArgs() ([]*Term, error) {
	var args []*Term

	// <args> ::= <term> <other args>
	term, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	args = append(args, term)

	// Parse <other args>
	otherArgs, err := p.parseOtherArgs()
	if err != nil {
		return nil, err
	}
	args = append(args, otherArgs...)

	return args, nil
}

// parseOtherArgs parses the <other args> non-terminal.
func (p *ParserImpl) parseOtherArgs() ([]*Term, error) {
	tok, err := p.peekToken()
	if err != nil {
		return nil, err
	}

	// <other args> ::= , <args> | ϵ
	if tok.typ != tokenComma {
		return nil, nil // Epsilon production
	}

	// Consume ','
	_, err = p.nextToken()
	if err != nil {
		return nil, err
	}

	// Parse <args>
	return p.parseArgs()
}

// getOrCreateTerm returns a cached term or creates a new one.
func (p *ParserImpl) getOrCreateTerm(typ TermType, literal string) *Term {
	key := strconv.Itoa(int(typ)) + ":" + literal
	if term, ok := p.usedTerms[key]; ok {
		return term
	}
	term := &Term{Typ: typ, Literal: literal}
	p.usedTerms[key] = term
	return term
}

// getOrCreateCompoundTerm returns a cached compound term or creates a new one.
func (p *ParserImpl) getOrCreateCompoundTerm(functor *Term, args []*Term) *Term {
	key := "compound:" + functor.Literal + "(" + TermSliceToString(args) + ")"
	if term, ok := p.usedTerms[key]; ok {
		return term
	}
	term := &Term{Typ: TermCompound, Functor: functor, Args: args}
	p.usedTerms[key] = term
	return term
}