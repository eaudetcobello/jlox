package ast

import jlox "github.com/eaudetcobello/jlox/pkg"

type Parser struct {
	tokens  []jlox.Token
	current int
}

func (p *Parser) ifMatchEat(types ...jlox.TokenType) bool {
	for _, t := range types {
		if p.currentTokenHasType(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) currentTokenHasType(t jlox.TokenType) bool {
	if p.tokens[p.current].TokenType == t {
		return true
	}

	return false
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == jlox.EOF
}

func (p *Parser) peek() jlox.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() jlox.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() jlox.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

//func (p *Parser) expression() Expr {
//	return p.equality()
//}

//func (p *Parser) equality() Expr {
//	expr := p.comparison()
//
//	for p.ifMatchEat(jlox.BANG_EQUAL, jlox.EQUAL_EQUAL) {
//		operator := p.previous()
//		right := p.comparison()
//		expr = Binary{
//			expr,
//			operator,
//			right,
//		}
//	}
//}
