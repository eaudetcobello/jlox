package jlox

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func RunFile(stdin io.Reader, stdout io.Writer, stderr io.Writer, fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if errs := Run(string(fileContent)); errs != nil {
		PrintError(stderr, errs)
	}

	return nil
}

func RunInteractive(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	r := bufio.NewReader(stdin)

	for {
		fmt.Fprint(stdout, "> ")

		s, _ := r.ReadString('\n')
		if s == "" {
			break
		}

		if errs := Run(s); errs != nil {
			PrintError(stderr, errs)

			continue
		}
	}

	return nil
}

func Run(s string) ErrorList {
	scanner := NewScanner(s)
	if err := scanner.scanTokens(); err != nil {
		return err
	}

	parser := Parser{
		Tokens:  scanner.Tokens,
		Current: 0,
		Errors:  ErrorList{},
	}

	var ast Expr

	var err ErrorList

	if ast, err = parser.parse(); err != nil {
		return err
	}

	astPrinter := AstPrinter{}
	fmt.Println(astPrinter.Print(ast))

	return nil
}
