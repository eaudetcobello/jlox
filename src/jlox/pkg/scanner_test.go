package jlox

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewScanner(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("abc123")
	g.Expect(scanner).NotTo(BeNil())

	g.Expect(scanner.Tokens).To(BeNil())
	g.Expect(scanner.Source).To(Equal("abc123"))
	g.Expect(scanner.Current).To(Equal(0))
	g.Expect(scanner.Line).To(Equal(1))
}

func TestScanErrorInPrompt(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("~~~")
	errs := scanner.scanTokens()

	g.Expect(errs).To(HaveLen(3))
	for _, err := range errs {
		g.Expect(err).To(HaveOccurred())
		g.Expect(err).To(MatchError("error: unexpected character"))
	}
}

func TestScanTokensInPrompt(t *testing.T) {
	g := NewWithT(t)

	for _, tc := range []struct {
		source         string
		ptrValue       *any
		expectedTokens []Token
	}{
		{
			source: "<>!==!.{=}",
			expectedTokens: []Token{
				{TokenType: LESS, Lexeme: "<", Literal: nil, Line: 1},
				{TokenType: GREATER, Lexeme: ">", Literal: nil, Line: 1},
				{TokenType: BANG_EQUAL, Lexeme: "!=", Literal: nil, Line: 1},
				{TokenType: EQUAL, Lexeme: "=", Literal: nil, Line: 1},
				{TokenType: BANG, Lexeme: "!", Literal: nil, Line: 1},
				{TokenType: DOT, Lexeme: ".", Literal: nil, Line: 1},
				{TokenType: LEFT_BRACE, Lexeme: "{", Literal: nil, Line: 1},
				{TokenType: EQUAL, Lexeme: "=", Literal: nil, Line: 1},
				{TokenType: RIGHT_BRACE, Lexeme: "}", Literal: nil, Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "(( )){}",
			expectedTokens: []Token{
				{TokenType: LEFT_PARENS, Lexeme: "(", Literal: nil, Line: 1},
				{TokenType: LEFT_PARENS, Lexeme: "(", Literal: nil, Line: 1},
				{TokenType: RIGHT_PARENS, Lexeme: ")", Literal: nil, Line: 1},
				{TokenType: RIGHT_PARENS, Lexeme: ")", Literal: nil, Line: 1},
				{TokenType: LEFT_BRACE, Lexeme: "{", Literal: nil, Line: 1},
				{TokenType: RIGHT_BRACE, Lexeme: "}", Literal: nil, Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "\"abc123\"",
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"abc123\"", Literal: "abc123", Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "123.4",
			expectedTokens: []Token{
				{TokenType: NUMBER, Lexeme: "123.4", Literal: 123.4, Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "\"123.4\"",
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"123.4\"", Literal: "123.4", Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "\"123.4\" //comment 123.4",
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"123.4\"", Literal: "123.4", Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "//comment",
			expectedTokens: []Token{
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "",
			expectedTokens: []Token{
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
	} {
		t.Run(fmt.Sprintf("scan %s", tc.source), func(t *testing.T) {
			scanner := NewScanner(tc.source)
			err := scanner.scanTokens()
			g.Expect(err).To(BeNil())
			g.Expect(scanner.Tokens).To(Equal(tc.expectedTokens))
		})
	}
}
