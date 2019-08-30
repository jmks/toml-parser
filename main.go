package main

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WHITESPACE

	// Literals
	KEY     // some_key, some.key, "key with spaces"
	STRING  // "some string"
	COMMENT // # comment til end of line
	INTEGER

	// Syntax
	ASSIGNMENT   // =
	BRAKET_OPEN  // [
	BRAKET_CLOSE // ]
	COMMA        // ,

	// Keywords
)

var eof = rune(0)

// Scanner is our TOML lexical scanner
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of the lexical scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	rune, _, err := s.r.ReadRune()

	if err != nil {
		return eof
	}

	return rune
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (Token, string) {
	rune := s.read()

	if isWhitespace(rune) {
		s.unread()
		return s.scanWhiteSpace()
	} else if isNumeric(rune) {
		s.unread()
		return s.scanNumeric()
	} else if isBareKeyLetter(rune) {
		s.unread()
		return s.scanKey()
	}

	switch rune {
	case eof:
		return EOF, ""
	case '[':
		return BRAKET_OPEN, "["
	case ']':
		return BRAKET_CLOSE, "]"
	case ',':
		return COMMA, ","
	case '#':
		s.unread()
		return COMMENT, s.scanToEndOfLine()
	case '"':
		return STRING, s.scanString()
	case '=':
		return ASSIGNMENT, "="
	}

	return ILLEGAL, string(rune)
}

func (s *Scanner) scanWhiteSpace() (Token, string) {
	var buffer bytes.Buffer

	buffer.WriteRune(s.read())

	for {
		rune := s.read()

		if rune == eof {
			break
		} else if !isWhitespace(rune) {
			s.unread()
			break
		} else {
			buffer.WriteRune(rune)
		}
	}

	return WHITESPACE, buffer.String()
}

// Note: does not support dotted keys
func (s *Scanner) scanKey() (Token, string) {
	var buffer bytes.Buffer

	buffer.WriteRune(s.read())

	for {
		rune := s.read()

		if rune == eof {
			break
		} else if !isBareKeyLetter(rune) {
			s.unread()
			break
		} else {
			buffer.WriteRune(rune)
		}
	}

	return KEY, buffer.String()
}

func (s *Scanner) scanNumeric() (Token, string) {
	var buffer bytes.Buffer

	for {
		rune := s.read()

		if rune == eof {
			break
		} else if !isDigit(rune) {
			s.unread()
			break
		} else {
			buffer.WriteRune(rune)
		}
	}

	return INTEGER, buffer.String()
}

// scanString is a bit naive right now
// It should also support escaped double quotes
func (s *Scanner) scanString() string {
	var buffer bytes.Buffer

	for {
		rune := s.read()

		if rune == eof {
			break
		} else if rune == '"' {
			break
		} else {
			buffer.WriteRune(rune)
		}
	}

	return buffer.String()
}

func (s *Scanner) scanToEndOfLine() string {
	rest, _ := s.r.ReadString('\n')

	return rest
}

func main() {}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

func isDashOrUnderscore(r rune) bool {
	return r == '-' || r == '_'
}

func isBareKeyLetter(r rune) bool {
	return isLetter(r) || isDigit(r) || isDashOrUnderscore(r)
}

func isNumeric(r rune) bool {
	return isDigit(r)
}
