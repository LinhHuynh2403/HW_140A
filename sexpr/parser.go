package sexpr

import (
	"errors"
	"math/big"
	"strings"
	"unicode"
)

var ErrParser = errors.New("parser error")

type Parser interface {
	Parse(string) (*SExpr, error)
}

func NewParser() Parser {
	return &sexprParser{}
}

type sexprParser struct{}

func (p *sexprParser) Parse(input string) (*SExpr, error) {
	tokens := tokenize(input)
	expr, remainingTokens, err := p.parseSExpr(tokens)
	if err != nil {
		return nil, err
	}
	if len(remainingTokens) > 0 {
		return nil, ErrParser
	}
	return expr, nil
}

func tokenize(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	for _, r := range input {
		if unicode.IsSpace(r) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else if r == '(' || r == ')' || r == '.' || r == '\'' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			currentToken.WriteRune(r)
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}
	return tokens
}

func (p *sexprParser) parseSExpr(tokens []string) (*SExpr, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, ErrParser
	}
	token := tokens[0]
	switch token {
	case "(":
		return p.parseList(tokens[1:])
	case "'":
		// Handle quote shorthand, e.g., 'x -> (QUOTE x)
		expr, remainingTokens, err := p.parseSExpr(tokens[1:])
		if err != nil {
			return nil, nil, err
		}
		return mkConsCell(mkSymbol("QUOTE"), mkConsCell(expr, mkNil())), remainingTokens, nil
	default:
		if is_number(token) {
			// Parse number
			num, _ := new(big.Int).SetString(token, 10)
			return mkNumber(num), tokens[1:], nil
		} else if isSymbol(token) {
			// Parse symbol
			return mkSymbol(token), tokens[1:], nil
		}
		return nil, nil, ErrParser
	}
}

func (p *sexprParser) parseList(tokens []string) (*SExpr, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, ErrParser
	}
	token := tokens[0]
	if token == ")" {
		// End of list
		return mkNil(), tokens[1:], nil
	}

	// Parse the first element of the list
	first, remainingTokens, err := p.parseSExpr(tokens)
	if err != nil {
		return nil, nil, err
	}

	// Check if the next token is a dot (for dotted lists)
	if len(remainingTokens) > 0 && remainingTokens[0] == "." {
		// Parse the second element of the dotted list
		second, remainingTokens, err := p.parseSExpr(remainingTokens[1:])
		if err != nil {
			return nil, nil, err
		}
		// Expect a closing parenthesis after the dotted list
		if len(remainingTokens) == 0 || remainingTokens[0] != ")" {
			return nil, nil, ErrParser
		}
		return mkConsCell(first, second), remainingTokens[1:], nil
	}

	// Parse the rest of the list
	rest, remainingTokens, err := p.parseList(remainingTokens)
	if err != nil {
		return nil, nil, err
	}
	return mkConsCell(first, rest), remainingTokens, nil
}

func is_number(s string) bool {
	_, ok := new(big.Int).SetString(s, 10)
	return ok
}

func isSymbol(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && r != '+' && r != '-' && r != '*' && r != '/' {
			return false
		}
	}
	return true
}
