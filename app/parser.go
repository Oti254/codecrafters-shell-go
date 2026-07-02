package main

import (
	"fmt"
	"strings"
)

type Redirection struct {
	Filename      string
	RedirectFound bool
	WorkingArgs   []string
}

// Handling redirections
func analyzeRedirection(args []string) (Redirection, error) {
	r := Redirection{"", false, args}
	// Analyze the command line arguments
	// TODO:
	// Support multiple redirections (e.g. "echo hello > a > b").
	// The parser currently stops after the first output redirection because
	// the current shell implementation supports only a single stdout redirect.
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			r.RedirectFound = true
			r.WorkingArgs = args[:i]
			if len(args) <= (i + 1) {
				return r, fmt.Errorf("shell: parse error near '\\n'\n")
			}
			r.Filename = args[i+1]
		}
	}
	return r, nil
}

// Handling single quotes
func parseCommand(command string) []string {
	var (
		words          []string
		current        strings.Builder
		inSingleQuotes bool
		inDoubleQuotes bool
	)

	for i := 0; i < len(command); i++ {
		char := command[i]
		switch {
		// Handling double quotes
		case char == '"' && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes

		// Switching the state
		case char == '\'' && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes

		// Handling backslash outside of quotes
		case char == '\\' && !inDoubleQuotes && !inSingleQuotes:
			current.WriteByte(command[i+1])
			i++

		// Handling backslash inside of double quotes
		// This should enable escaping of characters
		case char == '\\' && inDoubleQuotes:
			current.WriteByte(command[i+1])
			i++

		// Handling the spaces outside the quotes
		// When we encounter a space save the word
		// The space is treated as a separator not a character
		case char == ' ' && !inSingleQuotes && !inDoubleQuotes:
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}

		// Writing the characters by default to the current container housing characters
		// Before moving them to the string
		default:
			current.WriteByte(char)
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}
