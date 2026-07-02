package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// List of all builtin types
	builtIn := map[string]bool{
		"echo": true,
		"exit": true,
		"type": true,
		"pwd":  true,
		"cd":   true,
	}

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
		words := parseCommand(command)

		if len(words) == 0 {
			continue
		}

		// Getting the contents of the PATH variable
		pathEnv := os.Getenv("PATH")
		paths := filepath.SplitList(pathEnv)

		cmd := words[0]
		args := words[1:]

		info, err := analyzeRedirection(args)
		if err != nil {
			fmt.Println(err)
		}
		var writer io.Writer
		// Configuring where the child process is written to
		if info.RedirectFound {
			// Configuring the child process
			file, err := os.OpenFile(info.Filename, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			// Writes the child process directly to the file redirected
			writer = file
		} else {
			// Writes the child process directly to the terminal
			writer = os.Stdout
		}

		switch cmd {
		case "exit":
			return

		case "echo":
			handleEcho(writer, info.WorkingArgs)

		case "type":
			handleType(writer, words, builtIn, paths)

		case "pwd":
			handlePWD(writer)

		case "cd":
			handleCD(writer, cmd, args)

		default:
			handleProgram(writer, cmd, info.WorkingArgs)
		}
	}
}
