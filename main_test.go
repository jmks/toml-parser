package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	toml := "key = \"value\""

	want := []Token{BARE_KEY, WHITESPACE, ASSIGNMENT, WHITESPACE, STRING, EOF}
	got := Parse(toml)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %v, but got %v", want, got)
	}
}
