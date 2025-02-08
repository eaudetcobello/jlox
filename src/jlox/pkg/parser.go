package jlox

type Parser struct {
	Tokens  []Token
	Current int
	Errors  ErrorList
}

func (p *Parser) parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.ifMatchEat(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.ifMatchEat(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.ifMatchEat(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.ifMatchEat(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{
			expr,
			operator,
			right,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.ifMatchEat(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{
			operator,
			right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.ifMatchEat(FALSE) {
		return Literal{false}
	}

	if p.ifMatchEat(TRUE) {
		return Literal{true}
	}

	if p.ifMatchEat(NIL) {
		return Literal{nil}
	}

	if p.ifMatchEat(NUMBER, STRING) {
		return Literal{p.previous().Literal}
	}

	if p.ifMatchEat(LEFT_PARENS) {
		expr := p.expression()
		p.mustEat(RIGHT_PARENS)
		return Grouping{expr}
	}

	return nil
}

func (p *Parser) ifMatchEat(types ...TokenType) bool {
	for _, t := range types {
		if p.currentTokenHasType(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) mustEat(t TokenType) {
	if !p.ifMatchEat(t) {
		panic("Expected token not found")
	}
}

func (p *Parser) currentTokenHasType(t TokenType) bool {
	if p.Tokens[p.Current].TokenType == t {
		return true
	}

	return false
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

//func (p *Parser) expression() Expr {
//	return p.equality()
//}

//func (p *Parser) equality() Expr {
//	expr := p.comparison()
//
//	for p.ifMatchEat(BANG_EQUAL, EQUAL_EQUAL) {
//		operator := p.previous()
//		right := p.comparison()
//		expr = Binary{
//			expr,
//			operator,
//			right,
//		}
//	}
//}
