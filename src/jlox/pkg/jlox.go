package jlox

import (
	"bufio"
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
		if s == "" || s == "exit" {
			break
		}

		if err := Run(s); err != nil {
			return err
		}
	}

	return nil
}

func Run(s string) error {
	scanner := NewScanner(s)
	if err := scanner.scanTokens(); err != nil {
		return err
	}

	return nil
}
