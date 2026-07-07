package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
		words := parseCommand(command)

		if len(words) == 0 {
			continue
		}

		cmd := words[0]
		args := words[1:]

		info, err := analyzeRedirection(args)
		if err != nil {
			fmt.Println(err)
		}
		var w io.Writer
		// Configuring where the child process is written to
		if info.RedirectFound {
			// Configuring the child process
			file, err := os.OpenFile(info.Filename, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			// Writes the child process directly to the file redirected
			w = file
		} else {
			// Writes the child process directly to the terminal
			w = os.Stdout
		}

		handler, ok := builtinRegistry[cmd]
		if ok {
			handler(w, words)
		} else {
			handleProgram(w, cmd, info.WorkingArgs)
		}

		/**

		switch cmd {
		case "exit":
			return

		case "echo":
			handleEcho(w, info.WorkingArgs)

		case "type":
			handleType(w, words)

		case "pwd":
			handlePWD(w, words)

		case "cd":
			handleCD(w, cmd, args)

		default:
			handleProgram(w, cmd, info.WorkingArgs)
		}
		**/
	}
}
