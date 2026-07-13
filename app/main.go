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
		words, err := parseCommand(command)
		if err != nil {
			fmt.Println(err)
		}

		if words.Name == "" {
			continue
		}

		cmd := words.Name
		args := words.Args

		var stdout io.Writer = os.Stdout
		var stderr io.Writer = os.Stderr
		// Configuring where the child process is written to
		for _, redir := range words.Redirections {
			switch {
			case redir.FD == 1:
				// Configuring the child process
				file, err := os.OpenFile(redir.Filename, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()

				// Writes the child process directly to the file redirected
				stdout = file
			case redir.FD == 2:
				// Configuring the child process
				file, err := os.OpenFile(redir.Filename, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()

				// Writes the child process directly to the file redirected
				stderr = file

			}
		}
		/**
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
		**/
		input := append([]string{cmd}, args...)

		handler, ok := builtinRegistry[cmd]
		if ok {
			handler(stdout, input)
		} else {
			handleProgram(stdout, stderr, cmd, args)
		}
	}
}
