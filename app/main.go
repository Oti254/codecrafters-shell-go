package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Infinite loop for the REPL
	for {
		// Prints the prompt
		fmt.Print("$ ")

		// Reads the user input, stores it in a string
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input", err)
			os.Exit(1)
		}

		// Removes the newline at the end
		command = strings.TrimSpace(command)

		// Eliminates the spaces and places individual words in a list
		// Factoring in words that are inside quotes
		// This parses the commands from the command line
		words, err := parseCommand(command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error", err)
		}

		if words.Name == "" {
			continue
		}
		// Execute the command and handle any errors
		if err := executeCmd(words); err != nil {
			// Main decides how to display errors
			var exitError *exec.ExitError
			var execError *exec.Error

			if errors.As(err, &execError) {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", words.Name)
			} else if errors.As(err, &exitError) {
				// Program ran but exited with non-zero status
				// Don't print anything - the program's stderr output was already captured
				continue
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}
