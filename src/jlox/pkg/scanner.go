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
		TokenType: t,
		Lexeme:    scanner.Source[scanner.Start:scanner.Current],
		Literal:   &literal,
		Line:      scanner.Line,
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

	eofToken := Token{
		TokenType: EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      scanner.Line,
	}
	scanner.Tokens = append(scanner.Tokens, eofToken)

	for i, t := range scanner.Tokens {
		fmt.Printf("%d %s\n", i+1, t.String())
	}

	return nil
}

func (scanner *Scanner) consume() byte {
	char := scanner.Source[scanner.Current]
	scanner.Current++
	return char
}

func (scanner *Scanner) matchAndConsume(expected string) bool {
	if scanner.isAtEnd() {
		return false
	}

	if string(scanner.Source[scanner.Current]) != expected {
		return false
	}

	scanner.Current++

	return true
}

func (scanner *Scanner) peek() string {
	if scanner.isAtEnd() {
		return "\\0"
	}
	return string(scanner.Source[scanner.Current])
}

func (scanner *Scanner) string() *JLoxError {
	for scanner.peek() != "\"" && !scanner.isAtEnd() {
		if scanner.peek() == "\n" {
			scanner.Line++
		}
		scanner.consume()
	}

	if scanner.isAtEnd() {
		return NewJLoxError(fmt.Errorf("unterminated string"), false)
	}

	scanner.consume()

	strVal := scanner.Source[scanner.Start+1 : scanner.Current-1]
	scanner.addToken(STRING, strVal)

	return nil
}

func (scanner *Scanner) scanToken() *JLoxError {
	switch char := string(scanner.consume()); char {
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
	case "!":
		if scanner.matchAndConsume("=") {
			scanner.addTokenNoLiteral(BANG_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(BANG)
	case "=":
		if scanner.matchAndConsume("=") {
			scanner.addTokenNoLiteral(EQUAL_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(EQUAL)
	case "<":
		if scanner.matchAndConsume("=") {
			scanner.addTokenNoLiteral(LESS_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(LESS)
	case ">":
		if scanner.matchAndConsume("=") {
			scanner.addTokenNoLiteral(GREATER_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(GREATER)
	case string('"'):
		scanner.string()
	case "/":
		if scanner.matchAndConsume("/") {
			for scanner.peek() != "\n" && !scanner.isAtEnd() {
				scanner.consume()
			}
			break
		}
		scanner.addTokenNoLiteral(SLASH)
	case " ", "\r", "\t":
		break
	case "\n":
		scanner.Line++
	case "x":
		return NewJLoxError(fmt.Errorf("unexpected character: %s", char), true)
	default:
		return NewJLoxError(fmt.Errorf("unexpected character: %s", char), false)
	}

	return nil
}
