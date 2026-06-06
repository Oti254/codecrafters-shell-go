package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
		command = strings.TrimSpace(command)
		words := strings.Fields(command)

		// List of all builtin types
		builtIn := map[string]bool{
			"echo": true,
			"exit": true,
			"type": true,
		}
		if command == "exit" {
			break
		}

		// Getting the contents of the PATH variable
		pathEnv := os.Getenv("PATH")
		paths := filepath.SplitList(pathEnv)

		// Implementing echo my way
		if len(words) > 0 && words[0] == "echo" {
			fmt.Println(strings.Join(words[1:], " "))
			continue
		}

		// Type checking of commands
		if len(words) > 0 && words[0] == "type" {
			if len(words) < 2 {
				fmt.Println("type: missing argument")
				continue
			}

			if _, exists := builtIn[words[1]]; exists {
				fmt.Println(words[1] + " is a shell builtin")
				continue
			}

			// Checking if a filename exists
			found := false
			isExec := false

			for _, dir := range paths {
				path := filepath.Join(dir, words[1])
				fi, err := os.Stat(path)
				if err != nil {
					continue
				}

				// Check if file is a directory
				if fi.IsDir() {
					continue
				}

				// Checking if the file is executable
				mode := fi.Mode()
				isExec = mode&0111 != 0

				found = true

				// Prints the path if the file found is executable
				if found && isExec {
					fmt.Printf("%s is %s\n", words[1], path)
					break
				}
			}
			if !found {
				fmt.Println(words[1] + ": not found")
			}
			continue

		}
		// Printing error message
		fmt.Println(command + ": command not found")
	}
}
