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

	// List of all builtin types
	builtIn := map[string]bool{
		"echo": true,
		"exit": true,
		"type": true,
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
		command = strings.TrimSpace(command)
		words := strings.Fields(command)

		if len(words) == 0 {
			continue
		}

		// Getting the contents of the PATH variable
		pathEnv := os.Getenv("PATH")
		paths := filepath.SplitList(pathEnv)

		switch cmd := words[0]; cmd {
		case "exit":
			return

		case "echo":
			handleEcho(words)

		case "type":
			handleType(words, builtIn, paths)

		// Printing error message
		default:
			fmt.Printf("%s: not found", cmd)
		}
	}
}

// Implementing echo my way
func handleEcho(words []string) {
	fmt.Println(strings.Join(words[1:], " "))
}

// Type checking of commands
func handleType(words []string, builtIn map[string]bool, paths []string) {
	if len(words) < 2 {
		fmt.Println("type: missing argument")
	}

	if _, exists := builtIn[words[1]]; exists {
		fmt.Println(words[1] + " is a shell builtin")
		return
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
}
