package jlox

import (
	"bufio"
	"fmt"
	"io"
)

// RunInteractive runs each line read from stdin until stdin ends or fails to read.
func RunInteractive(stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	r := bufio.NewReader(stdin)

	for {
		fmt.Fprint(stdout, "> ")

		// ReadString returns what it read before EOF or a read error, so an
		// empty line with no newline means there is nothing left to run.
		s, _ := r.ReadString('\n')
		if s == "" {
			return
		}

		if errs := Run(stdout, stderr, s); errs != nil {
			PrintError(stderr, errs)
		}
	}
}

func Run(stdout io.Writer, stderr io.Writer, s string) error {
	scanner := NewScanner(s)
	if errs := scanner.scanTokens(); errs != nil {
		return errs
	}

	parser := Parser{
		Tokens: scanner.Tokens,
	}

	ast, errs := parser.parse()
	if errs != nil {
		return errs
	}

	result, err := evaluate(ast)
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, result)

	return nil
}
