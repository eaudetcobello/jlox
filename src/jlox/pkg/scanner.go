package jlox

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
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

	eofToken := Token{
		TokenType: EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      scanner.Line,
	}
	scanner.Tokens = append(scanner.Tokens, eofToken)

	if len(scanner.Errors) > 0 {
		return scanner.Errors
	}

	return nil
}

func (scanner *Scanner) eatCharacter() rune {
	char, width := utf8.DecodeRuneInString(scanner.Source[scanner.Current:])
	scanner.Current += width
	return char
}

func (scanner *Scanner) ifMatchEat(expected rune) bool {
	if scanner.isAtEnd() {
		return false
	}

	if rune(scanner.Source[scanner.Current]) != expected {
		return false
	}

	scanner.Current++

	return true
}

func (scanner *Scanner) Peek() rune {
	if scanner.isAtEnd() {
		return '\x00'
	}
	return rune(scanner.Source[scanner.Current])
}

func (scanner *Scanner) PeekNext() rune {
	if scanner.Current+1 >= len(scanner.Source) {
		return '\x00'
	}

	return rune(scanner.Source[scanner.Current+1])
}

func (scanner *Scanner) string() *Error {
	for scanner.Peek() != '"' && !scanner.isAtEnd() {
		if scanner.Peek() == '\n' {
			scanner.Line++
		}
		scanner.eatCharacter()
	}

	if scanner.isAtEnd() {
		return &Error{"unterminated string"}
	}

	scanner.eatCharacter()

	strVal := scanner.Source[scanner.Start+1 : scanner.Current-1]
	scanner.addToken(STRING, strVal)

	return nil
}

func (scanner *Scanner) number() *Error {
	for unicode.IsDigit(scanner.Peek()) {
		scanner.eatCharacter()
	}

	if scanner.Peek() == '.' && unicode.IsDigit(scanner.PeekNext()) {
		scanner.eatCharacter()

		for unicode.IsDigit(scanner.Peek()) {
			scanner.eatCharacter()
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
	switch char := scanner.eatCharacter(); char {
	case '(':
		scanner.addTokenNoLiteral(LEFT_PARENS)
	case ')':
		scanner.addTokenNoLiteral(RIGHT_PARENS)
	case '{':
		scanner.addTokenNoLiteral(LEFT_BRACE)
	case '}':
		scanner.addTokenNoLiteral(RIGHT_BRACE)
	case ',':
		scanner.addTokenNoLiteral(COMMA)
	case '.':
		scanner.addTokenNoLiteral(DOT)
	case '-':
		scanner.addTokenNoLiteral(MINUS)
	case '+':
		scanner.addTokenNoLiteral(PLUS)
	case ';':
		scanner.addTokenNoLiteral(SEMICOLON)
	case '*':
		scanner.addTokenNoLiteral(STAR)
	case '!':
		if scanner.ifMatchEat('=') {
			scanner.addTokenNoLiteral(BANG_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(BANG)
	case '=':
		if scanner.ifMatchEat('=') {
			scanner.addTokenNoLiteral(EQUAL_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(EQUAL)
	case '<':
		if scanner.ifMatchEat('=') {
			scanner.addTokenNoLiteral(LESS_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(LESS)
	case '>':
		if scanner.ifMatchEat('=') {
			scanner.addTokenNoLiteral(GREATER_EQUAL)
			break
		}
		scanner.addTokenNoLiteral(GREATER)
	case '"':
		return scanner.string()
	case '/':
		if scanner.ifMatchEat('/') {
			for scanner.Peek() != '\n' && !scanner.isAtEnd() {
				scanner.eatCharacter()
			}
			break
		}
		scanner.addTokenNoLiteral(SLASH)
	case ' ', '\r', '\t':
		break
	case '\n':
		scanner.Line++
	default:
		if unicode.IsDigit(char) {
			scanner.number()
		} else {
			return &Error{fmt.Sprintf("unexpected character %s", string(char))}
		}
	}

	return nil
}
