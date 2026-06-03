package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	// Prints the prompt
	fmt.Print("$ ")

	// Reads the user input, stores it in a string
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		errors.New("Unable to read user input")
	}

	// Printing error message
	fmt.Println(command[:len(command)-1] + ": command not found")
}
