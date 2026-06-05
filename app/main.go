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

		// Implementing echo
		if words[0] == "echo" {
			fmt.Print(words[1])
			for i, word := range words {
				if i == 0 || i == 1 {
					continue
				}
				fmt.Print(" " + word)
			}
			fmt.Println()
			continue
		}
		// Printing error message
		fmt.Println(command + ": command not found")
	}
}
