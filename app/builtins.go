package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// List of all builtin types
var builtIn = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

// Getting the current working directory
func handlePWD(w io.Writer) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprintln(w, path)
}

// Implementing echo my way
func handleEcho(w io.Writer, words []string) {
	fmt.Fprintln(w, strings.Join(words[:], " "))
}

func handleCD(w io.Writer, cmd string, args []string) {
	// Handling more than one argument
	if len(args) > 1 {
		fmt.Fprintf(w, "%s: too many arguments\n", cmd)
	}

	home := os.Getenv("HOME")
	var targetDir string

	switch {
	case len(args) == 0:
		targetDir = home
	case args[0] == "~":
		targetDir = home
	case strings.HasPrefix(args[0], "~/"):
		targetDir = filepath.Join(home, strings.TrimPrefix(args[0], "~/"))
	default:
		targetDir = args[0]
	}

	// Changing to the target directory
	err := os.Chdir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", cmd, targetDir)
	}
}

// Type checking of commands
func handleType(w io.Writer, words []string, builtIn map[string]bool, paths []string) {
	if len(words) < 2 {
		fmt.Fprintln(w, "type: missing argument")
		return
	}

	if _, exists := builtIn[words[1]]; exists {
		fmt.Fprintln(w, words[1]+" is a shell builtin")
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
			fmt.Fprintf(w, "%s is %s\n", words[1], path)
			break
		}
	}
	if !found {
		fmt.Fprintln(w, words[1]+": not found")
	}
}
