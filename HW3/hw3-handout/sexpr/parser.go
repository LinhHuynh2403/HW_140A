package sexpr

import (
	"errors"
)

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

/* Parsing table
            |		NUM/SYMBOl/QUOTE			|		 LPAR					|		RPAR			|		DOT					|      $        		||
<sexpr>     |	<sexpr> -> NUM/SYMBOl/QUOTE		|	<sexpr> -> LPAR <list> RPAR	|						|							|						||
            |									|								|						|							|						||
<list>		|	<list> -> <sexpr> <tail>		|	<list> -> <sexpr> <tail>	|	<list> -> epsilon	|							|	<list> -> epsilon	||
            |									|								|						|							|						||
<tail>      |	<tail> -> <sexpr> <tail>		|	<tail> -> <sexpr> <tail>	|	<tail> -> epsilon	|	<tail> -> DOT <sexpr>	|	<tail -> epsilon	||
*/

type Parser interface {
	Parse(string) (*SExpr, error)
}

type ParserImpl struct {
	lex       *lexer
	peekToken *token
}

func NewParser() Parser {
	return &ParserImpl{}
}

func (p *ParserImpl) Parse(input string) (*SExpr, error) {
	p.lex = newLexer(input)

	sexpr, err := p.parseSExpr()
	if err != nil {
		return nil, ErrParser
	}

	tok, err := p.nextToken()

	if tok.typ != tokenEOF {
		return nil, ErrParser
	}
	return sexpr, nil
}

func (p *ParserImpl) nextToken() (*token, error) {
	if p.peekToken != nil {
		tok := p.peekToken
		p.peekToken = nil
		return tok, nil
	}
	return p.lex.next()
}

func (p *ParserImpl) backToken(tok *token) {
	p.peekToken = tok
}

func (p *ParserImpl) parseSExpr() (*SExpr, error) {
    tok, err := p.nextToken()
    if err != nil {
    }

    switch tok.typ {
    case tokenNumber:
        return mkNumber(tok.num), nil
    case tokenSymbol:
        return mkSymbol(tok.literal), nil
    case tokenQuote:
        sexpr, err := p.parseSExpr()
        if err != nil {
            return nil, err
        }
        return mkConsCell(mkSymbol("QUOTE"), mkConsCell(sexpr, mkNil())), nil
    case tokenLpar:
        return p.parseList()
    default:
        return nil, ErrParser
    }
}

func (p *ParserImpl) parseList() (*SExpr, error) {
    tok, err := p.nextToken()
    if err != nil {
    }


    if tok.typ == tokenRpar {
        // Empty list
        return mkNil(), nil
    }

    p.backToken(tok)
    car, err := p.parseSExpr()
    if err != nil {
		
        return nil, err
    }

    cdr, err := p.parseTail()
    if err != nil {
        return nil, err
    }

    return mkConsCell(car, cdr), nil
}

func (p *ParserImpl) parseTail() (*SExpr, error) {
    tok, err := p.nextToken()
    if err != nil {
    }

    switch tok.typ {
    case tokenRpar:
        return mkNil(), nil
    case tokenDot:
        sexpr, err := p.parseSExpr()
        if err != nil {
            return nil, err
        }
        tok, err := p.nextToken()
        if err != nil {
        }
        if tok.typ != tokenRpar {
            // Invalid syntax: missing closing parenthesis after dot
            return nil, ErrParser
        }
        return sexpr, nil
    default:
        p.backToken(tok)
        car, err := p.parseSExpr()
        if err != nil {
            return nil, err
        }
        cdr, err := p.parseTail()
        if err != nil {
        }
        return mkConsCell(car, cdr), nil
    }
}