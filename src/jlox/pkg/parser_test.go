package jlox

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsAtEnd(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	g.Expect(p.isAtEnd()).To(BeFalse())

	p.advance()

	g.Expect(p.isAtEnd()).To(BeTrue())
}

func TestCurrentTokenHasType(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	g.Expect(p.currentTokenHasType(IF)).To(BeTrue())

	p.advance()

	g.Expect(p.currentTokenHasType(EOF)).To(BeTrue())
	g.Expect(p.currentTokenHasType(IF)).To(BeFalse())
}

func TestPeek(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	g.Expect(p.peek()).To(Equal(Token{
		TokenType: IF, Lexeme: "if", Literal: nil, Line: 1,
	}))

	p.advance()

	g.Expect(p.peek()).To(Equal(Token{
		TokenType: EOF, Lexeme: "", Literal: nil, Line: 1,
	}))
}

func TestAdvance(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	g.Expect(p.advance()).To(Equal(Token{
		TokenType: IF, Lexeme: "if", Literal: nil, Line: 1,
	}))
	g.Expect(p.Current).To(Equal(1))

	g.Expect(p.advance()).To(Equal(Token{
		TokenType: EOF, Lexeme: "", Literal: nil, Line: 1,
	}))
	g.Expect(p.Current).To(Equal(2))
}

func TestIfMatchEat(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	isMatch := p.ifMatchEat(AND)
	g.Expect(isMatch).To(BeFalse())
	g.Expect(p.Current).To(Equal(0))

	isMatch = p.ifMatchEat(OR, IF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.Current).To(Equal(1))

	isMatch = p.ifMatchEat(EOF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.Current).To(Equal(1))
}
