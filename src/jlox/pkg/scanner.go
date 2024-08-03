package jlox

import (
	"fmt"
)

type Scanner struct {
	tokens []Token
	source string

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		tokens:  nil,
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (scanner *Scanner) addToken(t TokenType, literal any) {
	scanner.tokens = append(scanner.tokens, Token{
		tokenType: t,
		lexeme:    scanner.source[scanner.start:scanner.current],
		literal:   literal,
		line:      scanner.line,
	})
}

func (scanner *Scanner) addTokenNoLiteral(t TokenType) {
	scanner.addToken(t, nil)
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) scanTokens() error {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		if err := scanner.scanToken(); err != nil {
			return err
		}
	}

	scanner.tokens = append(scanner.tokens, Token{EOF, "", nil, scanner.line})

	fmt.Println(scanner.tokens)

	return nil
}

func (scanner *Scanner) advance() byte {
	char := scanner.source[scanner.current]
	scanner.current++
	return char
}

func (scanner *Scanner) scanToken() error {
	switch char := string(scanner.advance()); char {
	case "(":
		scanner.addTokenNoLiteral(LEFT_PARENS)
	case ")":
		scanner.addTokenNoLiteral(RIGHT_PARENS)
	case "{":
		scanner.addTokenNoLiteral(LEFT_BRACE)
	case "}":
		scanner.addTokenNoLiteral(RIGHT_BRACE)
	case ",":
		scanner.addTokenNoLiteral(COMMA)
	case ".":
		scanner.addTokenNoLiteral(DOT)
	case "-":
		scanner.addTokenNoLiteral(MINUS)
	case "+":
		scanner.addTokenNoLiteral(PLUS)
	case ";":
		scanner.addTokenNoLiteral(SEMICOLON)
	case "*":
		scanner.addTokenNoLiteral(STAR)
	default:
		return fmt.Errorf("unexpected character: %s", char)
	}

	return nil
}
