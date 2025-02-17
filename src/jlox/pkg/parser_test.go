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

func TestCheck(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	g.Expect(p.check(IF)).To(BeTrue())

	p.advance()

	g.Expect(p.check(EOF)).To(BeTrue())
	g.Expect(p.check(IF)).To(BeFalse())
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

	g.Expect(p.advance()).To(Equal(Token{
		TokenType: EOF, Lexeme: "", Literal: nil, Line: 1,
	}))
	g.Expect(p.Current).To(Equal(2))
}

func TestMatch(t *testing.T) {
	g := NewWithT(t)

	p := Parser{
		Tokens: []Token{
			{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
			{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
		},
		Current: 0,
	}

	isMatch := p.match(AND)
	g.Expect(isMatch).To(BeFalse())
	g.Expect(p.Current).To(Equal(0))

	isMatch = p.match(OR, IF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.Current).To(Equal(1))

	isMatch = p.match(EOF)
	g.Expect(isMatch).To(BeTrue())
	g.Expect(p.Current).To(Equal(2))
}

func TestAst(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "2 + 2",
			expected: "(+ 2 2)",
		},
		{
			input:    "2 * 2 + 2",
			expected: "(+ (* 2 2) 2)",
		},
		{
			input:    "(2 * 2 + 2) ? 1 : 1",
			expected: "(conditional (group (+ (* 2 2) 2)) 1 1)",
		},
		{
			input:    "2 * 2 + 2 ? 1 : 1",
			expected: "(conditional (+ (* 2 2) 2) 1 1)",
		},
		{
			input:    "2 + 2,2,2,5 + 1",
			expected: "(comma (+ 2 2) 2 2 (+ 5 1))",
		},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)

		scanner.scanTokens()

		p := Parser{
			Tokens:  scanner.Tokens,
			Current: 0,
		}

		expr, errs := p.parse()

		astPrinter := AstPrinter{}

		g.Expect(astPrinter.Print(expr)).To(Equal(tt.expected))
		g.Expect(errs).To(BeNil())
	}
}
