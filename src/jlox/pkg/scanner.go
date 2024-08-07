package jlox

import (
	"fmt"
	"strconv"

	cmdutil "github.com/eaudetcobello/jlox/util"
)

type Scanner struct {
	Tokens []Token
	Source string
	Errors ErrorList

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

func (scanner *Scanner) addToken(t TokenType, literal interface{}) {
	if literal == nil {
		scanner.Tokens = append(scanner.Tokens, Token{
			TokenType: t,
			Lexeme:    scanner.Source[scanner.Start:scanner.Current],
			Literal:   nil,
			Line:      scanner.Line,
		})
	} else {
		scanner.Tokens = append(scanner.Tokens, Token{
			TokenType: t,
			Lexeme:    scanner.Source[scanner.Start:scanner.Current],
			Literal:   literal,
			Line:      scanner.Line,
		})
	}
}

func (scanner *Scanner) addTokenNoLiteral(t TokenType) {
	scanner.addToken(t, nil)
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)
}

func (scanner *Scanner) scanTokens() ErrorList {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		if err := scanner.scanToken(); err != nil {
			scanner.Errors.Add(err.Msg)
		}
	}

	if len(scanner.Errors) > 0 {
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

func (scanner *Scanner) Peek() string {
	if scanner.isAtEnd() {
		return "\\0"
	}
	return string(scanner.Source[scanner.Current])
}

func (scanner *Scanner) PeekNext() string {
	if scanner.Current+1 >= len(scanner.Source) {
		return "\\0"
	}

	return string(scanner.Source[scanner.Current+1])
}

func (scanner *Scanner) string() *Error {
	for scanner.Peek() != "\"" && !scanner.isAtEnd() {
		if scanner.Peek() == "\n" {
			scanner.Line++
		}
		scanner.consume()
	}

	if scanner.isAtEnd() {
		return &Error{"unterminated string"}
	}

	scanner.consume()

	strVal := scanner.Source[scanner.Start+1 : scanner.Current-1]
	scanner.addToken(STRING, strVal)

	return nil
}

func (scanner *Scanner) number() *Error {
	for cmdutil.IsDigit(scanner.Peek()) {
		scanner.consume()
	}

	if scanner.Peek() == "." && cmdutil.IsDigit(scanner.PeekNext()) {
		scanner.consume()

		for cmdutil.IsDigit(scanner.Peek()) {
			scanner.consume()
		}
	}

	v, err := strconv.ParseFloat(scanner.Source[scanner.Start:scanner.Current], 64)
	if err != nil {
		return &Error{"could not parse number"}
	}
	scanner.addToken(NUMBER, v)

	return nil
}

func (scanner *Scanner) scanToken() *Error {
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
			for scanner.Peek() != "\n" && !scanner.isAtEnd() {
				scanner.consume()
			}
			break
		}
		scanner.addTokenNoLiteral(SLASH)
	case " ", "\r", "\t":
		break
	case "\n":
		scanner.Line++
	default:
		if cmdutil.IsDigit(char) {
			scanner.number()
		} else {
			return &Error{"unexpected character"}
		}
	}

	return nil
}
