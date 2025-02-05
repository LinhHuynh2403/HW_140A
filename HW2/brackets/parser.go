package brackets

import "errors"

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

// Parser interface.
type Parser interface {
	Parse(string) (*Expr, error)
}

// Implement the Parser interface.
type ParserImpl struct {
	lex     *lexer
	peekTok *token
}

// Return a new instance of ParserImpl.
func NewParser() Parser {
	return &ParserImpl{}
}

// Helper function which returns the next token.
func (p *ParserImpl) nextToken() (*token, error) {
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

// Helper function which puts a token back as the next token.
func (p *ParserImpl) backToken(tok *token) {
	p.peekTok = tok
}

// Helper function to peek the next token.
func (p *ParserImpl) peekToken() (*token, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}

	p.backToken(tok)

	return tok, nil
}

// Implement the Parse method of the Parser interface for ParserImpl.
// Parse the non-terminal <start>.
// <start> ::= <bracket>
func (p *ParserImpl) Parse(input string) (*Expr, error) {
	p.lex = newLexer(input)

	// <start> -> <bracket>
	// Parse <bracket>.
	expr, err := p.bracketNT()
	if err != nil {
		return nil, ErrParser
	}

	// FOLLOW(<start>) = $
	// Check the next token is the endmarker $, there should be nothing left
	// after parsing <start>.
	if nextTok, err := p.nextToken(); err != nil || nextTok.typ != tokenEOF {
		return nil, ErrParser
	}

	return expr, nil
}

// Parse the non-terminal <bracket>.
// <bracket> ::= LBRKT <tail> RBRKT
func (p *ParserImpl) bracketNT() (*Expr, error) {

	// <bracket> -> LBRKT <tail> RBRKT
	// Consume LBRKT.
	tok, err := p.nextToken()
	if err != nil  || tok.typ != tokenLbracket {
		return nil, ErrParser
	}

	// Parse <tail>.
	expr, err := p.tailNT()
	if err != nil {
		return nil, ErrParser
	}

	// Consume RBRKT.
	tok, err = p.nextToken()
	if err != nil  || tok.typ != tokenRbracket {
		return nil, ErrParser
	}

	return expr, nil
}

// Parse the non-terminal <tail>.
// <tail> ::= <term> <tail> | \epsilon
func (p *ParserImpl) tailNT() (*Expr, error) {
	// There are two cases, so we peek the next token.
	tok, err := p.peekToken()
	if err != nil {
		return nil, ErrParser
	}

	switch tok.typ {

		// FIRST(<tail>)   =  {LBRKT, NUMBER, \epsilon}
		// FOLLOW(<tail>)  =  {RBRKT}
		// \epsilon in FIRST(<tail>) and RBRKT is in FOLLOW(<tail>).
		// <tail> -> \epsilon
		case tokenRbracket:
			return &Expr{Atom: nil, Bracket:[]*Expr{}}, nil

		// FIRST(<term>) = {LBRKT, NUMBER}
		// <tail> -> <term> <tail>
		case tokenLbracket, tokenNumber:
			// Parse <term>.
			term, err := p.termNT()
			if err != nil {
				return nil, ErrParser
			}

			// Parse <tail>.
			expr, err := p.tailNT()
			if err != nil {
				return nil, ErrParser
			}

			// Append the parsed <tail> to the parsed <term>.
			expr.Bracket = append([]*Expr{term}, expr.Bracket...)

			return expr, nil

		default:
			return nil, ErrParser

		}
}

// Parse the non-terminal <term>.
// <term> ::= NUMBER | <bracket>
func (p *ParserImpl) termNT() (*Expr, error) {
	// There are two cases, so we peek the next token.
	tok, err := p.peekToken()
	if err != nil {
		return nil, ErrParser
	}

	switch tok.typ {

	// <term> -> NUMBER
	case tokenNumber:
		// Consume this NUMBER token.
		tok, _ := p.nextToken()
		return &Expr{Atom: tok, Bracket:nil}, nil

	// LBRKT is in FIRST(<bracket>).
	// <term> -> <bracket>
	case tokenLbracket:
		return p.bracketNT()

	default:
		return nil, ErrParser

	}
}
