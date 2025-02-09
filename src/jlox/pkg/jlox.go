package jlox

import (
	"bufio"
	"fmt"
	"io"

	cmdutil "github.com/eaudetcobello/jlox/util"
)

type JLox struct{}

func RunFile(fileName string) error {
	fileContent, err := cmdutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	fmt.Println(string(fileContent))

	return nil
}

func RunPrompt(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	for {
		r := bufio.NewReader(stdin)

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

	for i, t := range scanner.Tokens {
		fmt.Printf("%d %s\n", i+1, t.String())
	}

	parser := Parser{
		Tokens:  scanner.Tokens,
		Current: 0,
		Errors:  ErrorList{},
	}

	var ast Expr

	var err ErrorList

	if ast, err = parser.parse(); err != nil {
		fmt.Println("Error parsing")
		return err
	}

	astPrinter := AstPrinter{}
	fmt.Println(astPrinter.Print(ast))

	return nil
}
