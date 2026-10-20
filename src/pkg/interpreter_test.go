package jlox

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

// parseSource scans and parses source, failing the test on any error.
func parseSource(g Gomega, source string) Expr {
	scanner := NewScanner(source)
	g.Expect(scanner.scanTokens()).To(BeNil())

	p := Parser{Tokens: scanner.Tokens}
	expr, errs := p.parse()
	g.Expect(errs).To(BeNil())

	return expr
}

func TestEvaluate(t *testing.T) {
	for _, tt := range []struct {
		source   string
		expected any
	}{
		{source: "1", expected: 1.0},
		{source: "2 * 3", expected: 6.0},
		{source: "1 + 2 * 3", expected: 7.0},
		{source: "(1 + 2) * 3", expected: 9.0},
		{source: "-(1 + 2)", expected: -3.0},
		{source: "!0", expected: false},
		{source: `"a" + "b"`, expected: "ab"},
		{source: "1 < 2", expected: true},
		{source: "1 == 1", expected: true},
		{source: `1 == "1"`, expected: false},
	} {
		t.Run(fmt.Sprintf("evaluate %q", tt.source), func(t *testing.T) {
			g := NewWithT(t)

			value, err := evaluate(parseSource(g, tt.source))

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(value).To(Equal(tt.expected))
		})
	}
}
