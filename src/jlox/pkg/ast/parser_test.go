package ast

import (
	"testing"

	jlox "github.com/eaudetcobello/jlox/pkg"
	. "github.com/onsi/gomega"
)

func TestIsAtEnd(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		tokens: []jlox.Token{
			{TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		current: 0,
	}

	g.Expect(p.isAtEnd()).To(BeFalse())

	p.advance()

	g.Expect(p.isAtEnd()).To(BeTrue())
}

func TestCurrentTokenHasType(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		tokens: []jlox.Token{
			{TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		current: 0,
	}

	g.Expect(p.currentTokenHasType(jlox.IF)).To(BeTrue())

	p.advance()

	g.Expect(p.currentTokenHasType(jlox.EOF)).To(BeTrue())
	g.Expect(p.currentTokenHasType(jlox.IF)).To(BeFalse())
}

func TestPeek(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		tokens: []jlox.Token{
			{TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		current: 0,
	}

	g.Expect(p.peek()).To(Equal(jlox.Token{
		TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1,
	}))

	p.advance()

	g.Expect(p.peek()).To(Equal(jlox.Token{
		TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1,
	}))
}

func TestAdvance(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		tokens: []jlox.Token{
			{TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		current: 0,
	}

	g.Expect(p.advance()).To(Equal(jlox.Token{
		TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1,
	}))
	g.Expect(p.current).To(Equal(1))

	g.Expect(p.advance()).To(Equal(jlox.Token{
		TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1,
	}))
	g.Expect(p.current).To(Equal(1))
}

func TestIfMatchEat(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		tokens: []jlox.Token{
			{TokenType: jlox.IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: jlox.EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		current: 0,
	}

	isMatch := p.ifMatchEat(jlox.AND)
	g.Expect(isMatch).To(BeFalse())
	g.Expect(p.current).To(Equal(0))

	isMatch = p.ifMatchEat(jlox.OR, jlox.IF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.current).To(Equal(1))

	isMatch = p.ifMatchEat(jlox.EOF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.current).To(Equal(1))
}
