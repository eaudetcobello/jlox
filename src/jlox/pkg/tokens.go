package jlox

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   *any
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("type: %s lexeme: %s literal: %v line: %d", t.tokenType, t.lexeme, "nil", t.line)
}

type TokenType int

const (
	LEFT_PARENS TokenType = iota
	RIGHT_PARENS
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	EOF
)

func (t TokenType) String() string {
	switch t {
	case LEFT_PARENS:
		return "LEFT_PARENS"
	case RIGHT_PARENS:
		return "RIGHT_PARENS"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
