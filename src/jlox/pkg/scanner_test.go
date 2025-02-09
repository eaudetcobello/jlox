package jlox

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

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

func TestEatCharacter(t *testing.T) {
	g := NewWithT(t)

	s := NewScanner("a😅1+ ⌘ 語\n\t")

	char := s.eatCharacter()
	g.Expect(char).To(Equal('a'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('😅'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('1'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('+'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal(' '))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('⌘'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal(' '))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('語'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('\n'))

	char = s.eatCharacter()
	g.Expect(char).To(Equal('\t'))
}

func TestScanTokensInPrompt(t *testing.T) {
	g := NewWithT(t)

	for _, tc := range []struct {
		source         string
		expectedTokens []Token
		expectedErrs   *ErrorList
	}{
		{
			source: "if and or ifandor",
			expectedTokens: []Token{
				{TokenType: IF, Lexeme: "if", Literal: nil, Line: 1},
				{TokenType: AND, Lexeme: "and", Literal: nil, Line: 1},
				{TokenType: OR, Lexeme: "or", Literal: nil, Line: 1},
				{TokenType: IDENTIFIER, Lexeme: "ifandor", Literal: nil, Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "123.4\"",
			expectedTokens: []Token{
				{TokenType: NUMBER, Lexeme: "123.4", Literal: 123.4, Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
			expectedErrs: &ErrorList{
				&Error{Msg: "unterminated string", Line: 1},
			},
		},
		{
			source: `"abc
def"`,
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"abc\ndef\"", Literal: "abc\ndef", Line: 2},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 2},
			},
		},
		{
			source: "<>!==!.{=}==<=>=+-;/*,",
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
				{TokenType: LESS_EQUAL, Lexeme: "<=", Literal: nil, Line: 1},
				{TokenType: GREATER_EQUAL, Lexeme: ">=", Literal: nil, Line: 1},
				{TokenType: PLUS, Lexeme: "+", Literal: nil, Line: 1},
				{TokenType: MINUS, Lexeme: "-", Literal: nil, Line: 1},
				{TokenType: SEMICOLON, Lexeme: ";", Literal: nil, Line: 1},
				{TokenType: SLASH, Lexeme: "/", Literal: nil, Line: 1},
				{TokenType: STAR, Lexeme: "*", Literal: nil, Line: 1},
				{TokenType: COMMA, Lexeme: ",", Literal: nil, Line: 1},
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
			source: "\"\"",
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"\"", Literal: "", Line: 1},
				{TokenType: EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			source: "\"    1 \"",
			expectedTokens: []Token{
				{TokenType: STRING, Lexeme: "\"    1 \"", Literal: "    1 ", Line: 1},
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
			s := NewScanner(tc.source)
			err := s.scanTokens()

			if tc.expectedErrs == nil {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).To(Equal(*tc.expectedErrs))
			}

			g.Expect(s.Tokens).To(Equal(tc.expectedTokens))
		})
	}
}
