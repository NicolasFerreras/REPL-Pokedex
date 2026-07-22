package repl

import (
	"strings"
)

// CleanInput parses and normalizes user input
func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)

	return words
}
