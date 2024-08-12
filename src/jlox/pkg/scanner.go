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

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	s.Tokens = append(s.Tokens, Token{
		TokenType: t,
		Lexeme:    s.Source[s.Start:s.Current],
		Literal:   literal,
		Line:      s.Line,
	})
}

func (s *Scanner) addTokenNoLiteral(t TokenType) {
	s.addToken(t, nil)
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanTokens() ErrorList {
	for !s.isAtEnd() {
		s.Start = s.Current
		if err := s.scanToken(); err != nil {
			s.Errors.Add(err.Msg)
		}
	}

	eofToken := Token{
		TokenType: EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.Line,
	}
	s.Tokens = append(s.Tokens, eofToken)

	if len(s.Errors) > 0 {
		return s.Errors
	}

	return nil
}

func (s *Scanner) eatCharacter() rune {
	char, width := utf8.DecodeRuneInString(s.Source[s.Current:])
	s.Current += width
	return char
}

func (s *Scanner) ifMatchEat(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.Source[s.Current]) != expected {
		return false
	}

	s.Current++

	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) peekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return '\x00'
	}

	return rune(s.Source[s.Current+1])
}

func (s *Scanner) string() *Error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++ // @We should store a string's first and last line
		}
		s.eatCharacter()
	}

	if s.isAtEnd() {
		return &Error{"unterminated string"}
	}

	s.eatCharacter()

	strVal := s.Source[s.Start+1 : s.Current-1]
	s.addToken(STRING, strVal)

	return nil
}

func (s *Scanner) number() *Error {
	for unicode.IsDigit(s.peek()) {
		s.eatCharacter()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.eatCharacter()

		for unicode.IsDigit(s.peek()) {
			s.eatCharacter()
		}
	}

	v, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	if err != nil {
		return &Error{"could not parse number"}
	}
	s.addToken(NUMBER, v)

	return nil
}

func (s *Scanner) identifier() *Error {
	char := s.peek()
	for unicode.IsLetter(char) || unicode.IsNumber(char) {
		s.eatCharacter()
		char = s.peek()
	}

	text := string(s.Source[s.Start:s.Current])

	token := keywords[text]

	if token == TokenType(0) { // @Hack?
		token = IDENTIFIER
	}
	s.addTokenNoLiteral(token)

	return nil
}

func (s *Scanner) scanToken() *Error {
	switch char := s.eatCharacter(); char {
	case '(':
		s.addTokenNoLiteral(LEFT_PARENS)
	case ')':
		s.addTokenNoLiteral(RIGHT_PARENS)
	case '{':
		s.addTokenNoLiteral(LEFT_BRACE)
	case '}':
		s.addTokenNoLiteral(RIGHT_BRACE)
	case ',':
		s.addTokenNoLiteral(COMMA)
	case '.':
		s.addTokenNoLiteral(DOT)
	case '-':
		s.addTokenNoLiteral(MINUS)
	case '+':
		s.addTokenNoLiteral(PLUS)
	case ';':
		s.addTokenNoLiteral(SEMICOLON)
	case '*':
		s.addTokenNoLiteral(STAR)
	case '!':
		if s.ifMatchEat('=') {
			s.addTokenNoLiteral(BANG_EQUAL)
			break
		}
		s.addTokenNoLiteral(BANG)
	case '=':
		if s.ifMatchEat('=') {
			s.addTokenNoLiteral(EQUAL_EQUAL)
			break
		}
		s.addTokenNoLiteral(EQUAL)
	case '<':
		if s.ifMatchEat('=') {
			s.addTokenNoLiteral(LESS_EQUAL)
			break
		}
		s.addTokenNoLiteral(LESS)
	case '>':
		if s.ifMatchEat('=') {
			s.addTokenNoLiteral(GREATER_EQUAL)
			break
		}
		s.addTokenNoLiteral(GREATER)
	case '"':
		return s.string()
	case '/':
		if s.ifMatchEat('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.eatCharacter()
			}
			break
		}
		s.addTokenNoLiteral(SLASH)
	case ' ', '\r', '\t':
		break
	case '\n':
		s.Line++
	default:
		if unicode.IsDigit(char) {
			s.number()
		} else if unicode.IsLetter(char) {
			s.identifier()
		} else {
			return &Error{fmt.Sprintf("unexpected character %s", string(char))}
		}
	}

	return nil
}
