package main

import (
	"fmt"
	"strings"
)

type Redirection struct {
	FD       int
	Operator string
	Filename string
}
type Command struct {
	Name         string
	Args         []string
	Redirections []Redirection
}

func parseCommand(input string) (Command, error) {
	words := tokenizer(input)
	cmd := Command{}

	// Checks if there are any words to parse
	if len(words) == 0 {
		return cmd, nil
	}

	for i := 0; i < len(words); i++ {
		switch {
		case i == 0:
			cmd.Name = words[i]
		case isRedirection(words[i]):
			switch {
			case len(words) <= (i + 1):
				return cmd, fmt.Errorf("shell: parse error near '\\n'\n")

			case words[i] == ">" || words[i] == "1>":
				cmd.Redirections = append(cmd.Redirections, Redirection{
					FD:       1,
					Operator: ">",
					Filename: words[i+1],
				})
			case words[i] == "2>":
				cmd.Redirections = append(cmd.Redirections, Redirection{
					FD:       2,
					Operator: ">",
					Filename: words[i+1],
				})
			}
			i += 1

		default:
			cmd.Args = append(cmd.Args, words[i])
		}
	}
	return cmd, nil
}

func tokenizer(input string) []string {
	var (
		words          []string
		current        strings.Builder
		inSingleQuotes bool
		inDoubleQuotes bool
	)

	for i := 0; i < len(input); i++ {
		char := input[i]
		switch {
		// Handling double quotes
		case char == '"' && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes

		// Switching the state
		case char == '\'' && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes

		// Handling backslash outside of quotes
		case char == '\\' && !inDoubleQuotes && !inSingleQuotes:
			current.WriteByte(input[i+1])
			i++

		// Handling backslash inside of double quotes
		// This should enable escaping of characters
		case char == '\\' && inDoubleQuotes:
			current.WriteByte(input[i+1])
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

var redirectChecker = map[string]struct{}{
	">":   {},
	"1>":  {},
	">>":  {},
	"1>>": {},
	"2>>": {},
	"2>":  {},
}

func isRedirection(token string) bool {
	_, exists := redirectChecker[token]
	return exists
}
