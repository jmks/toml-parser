package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		Description string
		Toml        string
		Tokens      []Token
	}{
		{
			Description: "A simple key-value",
			Toml:        "key = \"value\"",
			Tokens:      []Token{BARE_KEY, WHITESPACE, ASSIGNMENT, WHITESPACE, STRING, EOF},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Description, func(t *testing.T) {
			got := Parse(tc.Toml)

			if !reflect.DeepEqual(tc.Tokens, got) {
				t.Errorf("Wanted %v, but got %v", tc.Tokens, got)
			}
		})
	}
}
