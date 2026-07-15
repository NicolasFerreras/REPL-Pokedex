package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{

		{"hello world", []string{"hello", "world"}},
		{"  hello   world  ", []string{"hello", "world"}},
		{"one two three", []string{"one", "two", "three"}},
		{"Charmander Bulbasaur PIKACHU", []string{"charmander", "bulbasaur", "pikachu"}},
	}
	for _, tt := range tests {
		actual := cleanInput(tt.input)
		if len(actual) != len(tt.expected) {
			t.Errorf("cleanInput(%q) = %v; expected %v", tt.input, actual, tt.expected)
			continue
		}
		for i := range actual {
			if actual[i] != tt.expected[i] {
				t.Errorf("cleanInput(%q) = %v; expected %v", tt.input, actual, tt.expected)
				break
			}
		}
	}
}
