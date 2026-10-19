package jlox

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func RunFile(fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return Run(string(fileContent))
}

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

		if errs := Run(s); errs != nil {
			PrintError(stderr, errs)
		}
	}
}

func Run(s string) error {
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

	PrintAST(ast)

	return nil
}
