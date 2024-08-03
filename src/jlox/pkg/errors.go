package jlox

import "fmt"

type JLoxError struct {
	Err   error
	fatal bool
}

func (e *JLoxError) Error() string {
	return fmt.Sprintf("error: %v, irrecoverable: %v", e.Err, e.fatal)
}

func (e *JLoxError) IsFatal() bool {
	return e.fatal
}

func NewJLoxError(err error, fatal bool) *JLoxError {
	return &JLoxError{
		Err:   err,
		fatal: fatal,
	}
}
