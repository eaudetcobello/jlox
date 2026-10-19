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

func (s *Scanner) addToken(t TokenType, literal any) {
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
			s.Errors.Add(err)
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

func (s *Scanner) eatCharacter() byte {
	c := s.Source[s.Current]
	s.Current++

	return c
}

func (s *Scanner) ifMatchEat(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.Source[s.Current] != expected {
		return false
	}

	s.Current++

	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}

	return s.Source[s.Current]
}

func (s *Scanner) peekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return '\x00'
	}

	return s.Source[s.Current+1]
}

func (s *Scanner) string() error {
	startLine := s.Line

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++ // @We should store a string's first and last line
		}

		s.eatCharacter()
	}

	if s.isAtEnd() {
		return &Error{Msg: "unterminated string", Line: startLine}
	}

	s.eatCharacter()

	strVal := s.Source[s.Start+1 : s.Current-1]
	s.addToken(STRING, strVal)

	return nil
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.eatCharacter()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.eatCharacter()

		for isDigit(s.peek()) {
			s.eatCharacter()
		}
	}

	// The lexeme is only ASCII digits with at most one inner '.', so parsing cannot fail.
	v, _ := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	s.addToken(NUMBER, v)
}

func (s *Scanner) identifier() {
	char := s.peek()
	for isAlphaNumeric(char) {
		s.eatCharacter()
		char = s.peek()
	}

	text := string(s.Source[s.Start:s.Current])

	// is the identifier a known keyword
	if tokenType, ok := keywords[text]; ok {
		s.addTokenNoLiteral(tokenType)

		return
	}

	s.addTokenNoLiteral(IDENTIFIER)
}

func (s *Scanner) scanToken() error {
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
	case ':':
		s.addTokenNoLiteral(COLON)
	case '?':
		s.addTokenNoLiteral(QUESTION_MARK)
	case ' ', '\r', '\t':
		break
	case '\n':
		s.Line++
	default:
		if isDigit(char) {
			s.number()
		} else if isAlpha(char) {
			s.identifier()
		} else {
			// Failed to parse as ASCII, decode as UTF-8 to print correctly
			r, width := utf8.DecodeRuneInString(s.Source[s.Start:])
			s.Current = s.Start + width

			msg := fmt.Sprintf("unexpected character %q", r)
			if unicode.IsPrint(r) {
				msg = fmt.Sprintf("unexpected character '%c'", r)
			}
			return &Error{Msg: msg, Line: s.Line}
		}
	}

	return nil
}

// isAlpha reports whether c can start a Lox identifier.
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isDigit reports whether c is an ASCII digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
