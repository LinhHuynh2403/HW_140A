package sexpr

import "errors"

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <sexpr>       ::= <atom> | <pars> | QUOTE <sexpr>
// <atom>        ::= NUMBER | SYMBOL
// <pars>        ::= LPAR <dotted_list> RPAR | LPAR <proper_list> RPAR
// <dotted_list> ::= <proper_list> <sexpr> DOT <sexpr>
// <proper_list> ::= <sexpr> <proper_list> | \epsilon
//

// <sexpr>		 ::= NUMBER | SYMBOL | QUOTE <sexpr> | LPAR <list> RPAR
// <list>		 ::= <sexpr> <tail> | epsilon
// <tail>		 ::= epsilon | <sexpr> <tail> | DOT <sexpr>

//              |            FIRST            				|          FOLLOW          					||
//              |                             				|                          					||
// <sexpr>      | NUMBER, SYMBOL, QUOTE, LPAR           	| $, DOT, NUMBER, SYMBOL, QUOTE, LPAR, RPAR ||
//              |                             				|                          					||
// <list>       | NUMBER, SYMBOL, QUOTE, LPAR, epsilon  	| $, DOT, NUMBER, SYMBOL, QUOTE, LPAR, RPAR ||
//              |                	             			|                          					||
// <tail>       | DOT, NUMBER, SYMBOL, QUOTE, LPAR, epsilon | $, DOT, NUMBER, SYMBOL, QUOTE, LPAR, RPAR ||
//              |                             				|                          					||

type Parser interface {
	Parse(string) (*SExpr, error)
}

type ParserImpl struct {
	lex			*lexer
	peekToken 	*Token
	s_expr	 	[]*SExpr
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

func NewParser() Parser {
	// panic("TODO: implement NewParser")
	return &ParserImpl{
		lex: 		nil,
		peekToken: 	nil,
		s_expr: 	nil,
	}
}


