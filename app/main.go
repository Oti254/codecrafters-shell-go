package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/internal/terminal"
)

// Configuring tab completer choices
var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
)

func main() {
	// Replaces the bufio.NewScanner
	// Initializing the readline instance
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		log.Fatalf("Error initializing readline: %v", err)
	}
	defer l.Close()

	// Infinite loop for the REPL
	for {
		/**
		line, err := l.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			} else if err == io.EOF {
				break
			}
		}

		// Removes the newline at the end
		line = strings.TrimSpace(line)
		**/
		line, err := terminal.ReadCommand(l)
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			} else if err == io.EOF {
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
