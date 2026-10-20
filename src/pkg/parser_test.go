package jlox

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAstSuccessful(t *testing.T) {
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
		{
			input:    "1 ? 2 : 3, 4",
			expected: "(comma (conditional 1 2 3) 4)",
		},
		{
			input:    "1 + (2 * 3)",
			expected: "(+ 1 (group (* 2 3)))",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("parse %q", tt.input), func(t *testing.T) {
			g := NewWithT(t)

			scanner := NewScanner(tt.input)
			g.Expect(scanner.scanTokens()).To(BeNil())

			p := Parser{Tokens: scanner.Tokens}
			expr, errs := p.parse()

			g.Expect(errs).To(BeNil())
			g.Expect(PrintAST(expr)).To(Equal(tt.expected))
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, tt := range []struct {
		input   string
		message string
	}{
		{input: "2 . 2", message: "Expect end of expression."},
		{input: "1 2", message: "Expect end of expression."},
		{input: "1 +", message: "Expect expression."},
		{input: "(1", message: "Expected RIGHT_PARENS"},
	} {
		t.Run(fmt.Sprintf("parse %q", tt.input), func(t *testing.T) {
			g := NewWithT(t)

			scanner := NewScanner(tt.input)
			g.Expect(scanner.scanTokens()).To(BeNil())

			p := Parser{Tokens: scanner.Tokens}
			expr, errs := p.parse()

			g.Expect(expr).To(BeNil())
			g.Expect(errs).To(MatchError(ContainSubstring(tt.message)))
		})
	}
}
