package jlox

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	cmdutil "github.com/eaudetcobello/jlox/util"
)

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
		fmt.Fprint(stdout, "> ")
		r := bufio.NewReader(stdin)
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
	return nil
}
