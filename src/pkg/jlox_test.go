package jlox

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestRunInteractiveProcessesMultipleLines(t *testing.T) {
	g := NewWithT(t)

	input := "1 + 2\n3 + 4\n"
	stdin := strings.NewReader(input)
	var stdout, stderr bytes.Buffer

	RunInteractive(stdin, &stdout, &stderr)

	output := stdout.String()

	// With 2 input lines, we should see 3 prompts:
	// prompt 1 -> reads "1 + 2", prompt 2 -> reads "3 + 4", prompt 3 -> reads EOF and exits
	g.Expect(strings.Count(output, "> ")).To(Equal(3))
	g.Expect(stderr.String()).To(BeEmpty())
}
