package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/internal/terminal"
)

func main() {
	// Replaces the bufio.NewScanner
	// Initializing the readline instance
	term, err := terminal.New()
	if err != nil {
		log.Fatalf("Error initializing readline: %v", err)
	}
	defer term.Close()

	// Infinite loop for the REPL
	for {
		line, err := term.ReadCommand()
		if err != nil {
			if err == terminal.ErrInterrupt {
				continue
			} else if err == terminal.ErrEOF {
				break
			}
		}

		// This parses the commands from the command line
		words, err := parseCommand(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
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
