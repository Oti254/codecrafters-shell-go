package main

import (
	"bufio"
	"fmt"
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
		command = strings.TrimSpace(command)
		words := strings.Fields(command)

		if command == "exit" {
			break
		}

		// Implementing echo my way
		if len(words) > 0 && words[0] == "echo" {
			fmt.Println(strings.Join(words[1:], " "))
			continue
		}
		// Printing error message
		fmt.Println(command + ": command not found")
	}
}
