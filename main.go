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
	BARE_KEY   // some_key
	DOTTED_KEY // some.key
	STRING

	// Syntax
	ASSIGNMENT

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
	} else if isBareKeyLetter(rune) {
		s.unread()
		return s.scanBareKey()
	}

	switch rune {
	case eof:
		return EOF, ""
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

func (s *Scanner) scanBareKey() (Token, string) {
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

	return BARE_KEY, buffer.String()
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
		} else if !isBareKeyLetter(rune) {
			break
		} else {
			buffer.WriteRune(rune)
		}
	}

	return buffer.String()
}

func main() {}

func Parse(toml string) []Token {
	tokens := []Token{}

	scanner := NewScanner(bytes.NewBufferString(toml))

	for {
		token, _ := scanner.Scan()

		tokens = append(tokens, token)

		if token == EOF {
			break
		}
	}

	return tokens
}

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
