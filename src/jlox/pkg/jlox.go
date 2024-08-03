package jlox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	cmdutil "github.com/eaudetcobello/jlox/util"
)

type JLox struct {
}

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
		fmt.Fprint(stdout, "jlox> ")
		s, _ := r.ReadString('\n')
		s = strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s), " ", ""), "\n", "")
		if s == "" {
			break
		}

		if errs := Run(s); errs != nil {
			fmt.Println(errs)

			for _, err := range errs {
				if err.IsFatal() {
					return errors.New("encountered fatal error, exiting")
				}
			}

			continue
		}
	}

	return nil
}

func Run(s string) []*JLoxError {
	scanner := NewScanner(s)
	if err := scanner.scanTokens(); err != nil {
		return err
	}

	return nil
}
