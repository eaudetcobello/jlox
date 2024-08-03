package jlox

import (
	"fmt"
)

type Scanner struct {
	Tokens []Token
	Source string
	Errors []*JLoxError

	Start   int
	Current int
	Line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Tokens:  nil,
		Source:  source,
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

func (scanner *Scanner) addToken(t TokenType, literal any) {
	scanner.Tokens = append(scanner.Tokens, Token{
		tokenType: t,
		lexeme:    scanner.Source[scanner.Start:scanner.Current],
		literal:   literal,
		line:      scanner.Line,
	})
}

func (scanner *Scanner) addTokenNoLiteral(t TokenType) {
	scanner.addToken(t, nil)
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)
}

func (scanner *Scanner) scanTokens() []*JLoxError {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		if err := scanner.scanToken(); err != nil {
			scanner.Errors = append(scanner.Errors, err)
		}
	}

	if scanner.Errors != nil {
		return scanner.Errors
	}

	scanner.Tokens = append(scanner.Tokens, Token{EOF, "", nil, scanner.Line})

	fmt.Println(scanner.Tokens)

	return nil
}

func (scanner *Scanner) advance() byte {
	char := scanner.Source[scanner.Current]
	scanner.Current++
	return char
}

func (scanner *Scanner) scanToken() *JLoxError {
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
	case "x":
		return NewJLoxError(fmt.Errorf("unexpected character: %s", char), true)
	default:
		return NewJLoxError(fmt.Errorf("unexpected character: %s", char), false)
	}

	return nil
}
