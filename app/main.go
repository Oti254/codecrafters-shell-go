package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		words := strings.Fields(command)

		if len(words) == 0 {
			continue
		}

		// Getting the contents of the PATH variable
		pathEnv := os.Getenv("PATH")
		paths := filepath.SplitList(pathEnv)

		cmd := words[0]
		args := words[1:]

		switch cmd {
		case "exit":
			return

		case "echo":
			handleEcho(words)

		case "type":
			handleType(words, builtIn, paths)

		case "pwd":
			handlePWD()

		case "cd":
			handleCD(cmd, args)

		default:
			handleProgram(cmd, args)
		}
	}
}

func handleCD(cmd string, args []string) {
	// Handling more than one argument
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "%s: too many arguments\n", cmd)
	}

	var targetDir string

	// Handling no argument or ~
	if len(args) == 0 || args[0] == "~" {
		targetDir = os.Getenv("HOME")
	} else {
		targetDir = args[0]
	}

	// Changing to the target directory
	err := os.Chdir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", cmd, targetDir)
	}
}

// Getting the current working directory
func handlePWD() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(path)
}

// Implementing echo my way
func handleEcho(words []string) {
	fmt.Println(strings.Join(words[1:], " "))
}

// Running a program
func handleProgram(cmd string, args []string) {
	program := exec.Command(cmd, args...)

	// Writes the child process directly to the terminal
	program.Stdin = os.Stdin
	program.Stdout = os.Stdout
	program.Stderr = os.Stderr

	// Automatically checks if the program is in $PATH
	err := program.Run()
	if err != nil {
		fmt.Printf("%s: not found\n", cmd)
	}
}

// Type checking of commands
func handleType(words []string, builtIn map[string]bool, paths []string) {
	if len(words) < 2 {
		fmt.Println("type: missing argument")
		return
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
