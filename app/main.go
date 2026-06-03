package main

import (
	"bufio"
	"fmt"
	"os"
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

		// Printing error message
		fmt.Println(command[:len(command)-1] + ": command not found")
	}
}
