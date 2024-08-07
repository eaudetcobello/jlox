package jlox

import (
	"fmt"
	"io"
)

// From :
// https://cs.opensource.google/go/go/+/refs/tags/go1.22.6:src/go/scanner/errors.go

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %v", e.Msg)
}

type ErrorList []*Error

func (p *ErrorList) Add(msg string) {
	*p = append(*p, &Error{msg})
}

func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
