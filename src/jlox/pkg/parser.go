package jlox

type Parser struct {
	Tokens  []Token
	Current int
	Errors  ErrorList
}

func (p *Parser) parse() (Expr, ErrorList) {
	expr := p.expression()

	if len(p.Errors) > 0 {
		return nil, p.Errors
	}

	return expr, nil
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
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

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
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

	for p.match(MINUS, PLUS) {
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

	for p.match(SLASH, STAR) {
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
	if p.match(BANG, MINUS) {
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
	if p.match(FALSE) {
		return Literal{false}
	}

	if p.match(TRUE) {
		return Literal{true}
	}

	if p.match(NIL) {
		return Literal{nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().Literal}
	}

	if p.match(LEFT_PARENS) {
		expr := p.expression()

		if err := p.consume(RIGHT_PARENS); err != nil {
			p.Errors.Add(err)
			p.synchronize()
			return nil
		}

		return Grouping{expr}
	}

	return nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(t TokenType) error {
	if !p.match(t) {
		err := &Error{"Expected token of type " + t.String()}
		p.Errors.Add(err)
		p.synchronize()
		return err
	}
	return nil
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == t
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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
