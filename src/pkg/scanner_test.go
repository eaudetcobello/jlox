package jlox

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestScanErrorsInPrompt(t *testing.T) {
	g := NewWithT(t)

	scanner := NewScanner("125.4 ~~~ {}")
	errs := scanner.scanTokens()

	g.Expect(errs).To(HaveLen(3))
	g.Expect(errs).To(MatchError("[line 1] error: unexpected character '~' (and 2 more errors)"))

	for _, err := range errs {
		g.Expect(err).To(HaveOccurred())
		g.Expect(err).To(MatchError("[line 1] error: unexpected character '~'"))
	}

	g.Expect(scanner.Tokens).To(HaveLen(4))
}

func TestEatCharacter(t *testing.T) {
	g := NewWithT(t)

	s := NewScanner("a1+")

	g.Expect(s.eatCharacter()).To(Equal(byte('a')))
	g.Expect(s.eatCharacter()).To(Equal(byte('1')))
	g.Expect(s.eatCharacter()).To(Equal(byte('+')))
	g.Expect(s.isAtEnd()).To(BeTrue())
}

func TestScanNonASCIIOutsideString(t *testing.T) {
	g := NewWithT(t)

	s := NewScanner("a😅")
	errs := s.scanTokens()

	g.Expect(errs).To(HaveLen(1))
	g.Expect(errs[0]).To(MatchError("[line 1] error: unexpected character '😅'"))
	g.Expect(s.Tokens).To(HaveLen(2))
}

// tokenStrings renders tokens as "TYPE lexeme" so tests can ignore literal and line.
func tokenStrings(tokens []Token) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = strings.TrimSpace(t.TokenType.String() + " " + t.Lexeme)
	}

	return out
}

func TestScanTokensInPrompt(t *testing.T) {
	for _, tc := range []struct {
		source   string
		want     []string
		wantErrs ErrorList
	}{
		{
			source: "if and or ifandor",
			want:   []string{"IF if", "AND and", "OR or", "IDENTIFIER ifandor", "EOF"},
		},
		{
			source:   "123.4\"",
			want:     []string{"NUMBER 123.4", "EOF"},
			wantErrs: ErrorList{&Error{Msg: "unterminated string", Line: 1}},
		},
		{
			source: "<>!==!.{=}==<=>=+-;/*,",
			want: []string{
				"LESS <", "GREATER >", "BANG_EQUAL !=", "EQUAL =", "BANG !", "DOT .",
				"LEFT_BRACE {", "EQUAL =", "RIGHT_BRACE }", "EQUAL_EQUAL ==", "LESS_EQUAL <=",
				"GREATER_EQUAL >=", "PLUS +", "MINUS -", "SEMICOLON ;", "SLASH /", "STAR *",
				"COMMA ,", "EOF",
			},
		},
		{
			source: "(( )){}",
			want: []string{
				"LEFT_PARENS (", "LEFT_PARENS (", "RIGHT_PARENS )", "RIGHT_PARENS )",
				"LEFT_BRACE {", "RIGHT_BRACE }", "EOF",
			},
		},
		{
			source: `"123.4" //comment 123.4`,
			want:   []string{`STRING "123.4"`, "EOF"},
		},
		{
			source: "//comment",
			want:   []string{"EOF"},
		},
		{
			source: "",
			want:   []string{"EOF"},
		},
	} {
		t.Run(fmt.Sprintf("scan %q", tc.source), func(t *testing.T) {
			g := NewWithT(t)

			s := NewScanner(tc.source)
			errs := s.scanTokens()

			g.Expect(errs).To(Equal(tc.wantErrs))
			g.Expect(tokenStrings(s.Tokens)).To(Equal(tc.want))
		})
	}
}

func TestScanLiterals(t *testing.T) {
	for _, tc := range []struct {
		source string
		want   any
	}{
		{source: "123.4", want: 123.4},
		{source: `"abc123"`, want: "abc123"},
		{source: `"123.4"`, want: "123.4"},
		{source: `""`, want: ""},
		{source: `"    1 "`, want: "    1 "},
	} {
		t.Run(fmt.Sprintf("scan %q", tc.source), func(t *testing.T) {
			g := NewWithT(t)

			s := NewScanner(tc.source)

			g.Expect(s.scanTokens()).To(BeNil())
			g.Expect(s.Tokens).To(HaveLen(2))
			g.Expect(s.Tokens[0].Literal).To(Equal(tc.want))
		})
	}
}

func TestScanLines(t *testing.T) {
	for _, tc := range []struct {
		source string
		want   []int
	}{
		{source: "\"abc\ndef\"", want: []int{2, 2}},
		{source: "1\n2", want: []int{1, 2, 2}},
	} {
		t.Run(fmt.Sprintf("scan %q", tc.source), func(t *testing.T) {
			g := NewWithT(t)

			s := NewScanner(tc.source)
			g.Expect(s.scanTokens()).To(BeNil())

			lines := make([]int, len(s.Tokens))
			for i, tok := range s.Tokens {
				lines[i] = tok.Line
			}

			g.Expect(lines).To(Equal(tc.want))
		})
	}
}
