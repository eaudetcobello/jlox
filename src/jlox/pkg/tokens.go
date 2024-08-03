package jlox

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
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
