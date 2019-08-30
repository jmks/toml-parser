package main

import (
	"bytes"
	"fmt"
	"testing"
)

type TokenValue struct {
	Token
	Value string
}

func TestScanner(t *testing.T) {
	testcases := []struct {
		Toml        string
		TokenValues []TokenValue
	}{
		{
			Toml: "key = \"value\"",
			TokenValues: []TokenValue{
				TokenValue{KEY, "key"},
				TokenValue{WHITESPACE, " "},
				TokenValue{ASSIGNMENT, "="},
				TokenValue{WHITESPACE, " "},
				TokenValue{STRING, "value"},
			},
		},
		{
			Toml: "\"Quoted key\" = [1,2,3]  # Comments are ignored",
			TokenValues: []TokenValue{
				TokenValue{STRING, "Quoted key"},
				TokenValue{WHITESPACE, " "},
				TokenValue{ASSIGNMENT, "="},
				TokenValue{WHITESPACE, " "},
				TokenValue{BRAKET_OPEN, "["},
				TokenValue{INTEGER, "1"},
				TokenValue{COMMA, ","},
				TokenValue{INTEGER, "2"},
				TokenValue{COMMA, ","},
				TokenValue{INTEGER, "3"},
				TokenValue{BRAKET_CLOSE, "]"},
				TokenValue{WHITESPACE, "  "},
				TokenValue{COMMENT, "# Comments are ignored"},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("'%s' should parse to %v", tc.Toml, tc.TokenValues), func(t *testing.T) {
			scanner := NewScanner(bytes.NewBufferString(tc.Toml))

			for _, tv := range tc.TokenValues {
				token, value := scanner.Scan()

				if token != tv.Token {
					t.Errorf("Expected token %d, got %d", tv.Token, token)
				}

				if value != tv.Value {
					t.Errorf("Expected value '%s', got '%s'", tv.Value, value)
				}
			}
		})
	}
}
