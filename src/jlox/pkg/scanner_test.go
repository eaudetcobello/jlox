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

func TestScanErrorsInPrompt(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("125.4 ~~~ {}")
	errs := scanner.scanTokens()

	g.Expect(errs).To(HaveLen(3))
	g.Expect(errs).To(MatchError("error: unexpected character ~ (and 2 more errors)"))
	for _, err := range errs {
		g.Expect(err).To(HaveOccurred())
		g.Expect(err).To(MatchError("error: unexpected character ~"))
	}

	g.Expect(scanner.Tokens).To(HaveLen(4))
}

func TestUnterminatedString(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("123.4\"")
	errs := scanner.scanTokens()

	g.Expect(errs).To(HaveLen(1))
	g.Expect(errs[0]).To(MatchError("error: unterminated string"))
	g.Expect(scanner.Tokens[0].TokenType).To(Equal(NUMBER))
	g.Expect(scanner.Tokens[1].TokenType).To(Equal(EOF))
	g.Expect(scanner.Tokens).To(HaveLen(2))
}

func TestEatCharacter(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("a😅1+ ⌘ 語\n\t")

	char := scanner.eatCharacter()
	g.Expect(char).To(Equal('a'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('😅'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('1'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('+'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal(' '))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('⌘'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal(' '))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('語'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('\n'))

	char = scanner.eatCharacter()
	g.Expect(char).To(Equal('\t'))
}

func TestScanTokensInPrompt(t *testing.T) {
	g := NewWithT(t)

	for _, tc := range []struct {
		source         string
		ptrValue       *any
		expectedTokens []Token
	}{
		{
			source: "<>!==!.{=}==",
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
				{TokenType: EQUAL_EQUAL, Lexeme: "==", Literal: nil, Line: 1},
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
