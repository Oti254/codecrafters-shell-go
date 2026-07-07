package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Separating the checking of the existence of a command
// with the handler responsible for executing the command
// Implementing a set of builtin command names
var builtinChecker = map[string]struct{}{
	"cd":   {},
	"echo": {},
	"exit": {},
	"pwd":  {},
	"type": {},
}

type BuiltinHandler func(io.Writer, []string)

// Checking what is responsible for handling commands
// Instead of does a command exist

var builtinRegistry = map[string]BuiltinHandler{
	"cd":   handleCD,
	"echo": handleEcho,
	"exit": handleExit,
	"pwd":  handlePWD,
	"type": handleType,
}

/**
func getBuiltinRegistry() map[string]BuiltinHandler {
	return map[string]BuiltinHandler{
		"cd":   handleCD,
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePWD,
		"exit": handleExit,
	}
}
**/

func inRegistry(cmd string) bool {
	if _, exists := builtinChecker[cmd]; exists {
		return exists
	}
	return false
}

/**
// List of all builtin types
var builtIn = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}
**/

// Getting the contents of the PATH variable
var pathEnv = os.Getenv("PATH")
var paths = filepath.SplitList(pathEnv)

// Getting the current working directory
func handlePWD(w io.Writer, words []string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(w, path)
}

// Implementing echo my way
func handleEcho(w io.Writer, words []string) {
	fmt.Fprintln(w, strings.Join(words[1:], " "))
}

func handleCD(w io.Writer, words []string) {
	cmd := words[0]
	args := words[1:]
	// Handling more than one argument
	if len(args) > 1 {
		fmt.Fprintf(w, "%s: too many arguments\n", cmd)
		return
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
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", targetDir)
	}
}

func handleExit(w io.Writer, words []string) {
	os.Exit(0)
}

// Type checking of commands
func handleType(w io.Writer, words []string) {
	// Getting the contents of the PATH variable
	pathEnv = os.Getenv("PATH")
	paths = filepath.SplitList(pathEnv)
	if len(words) < 2 {
		fmt.Fprintln(w, "type: missing argument")
		return
	}

	cmd := words[1]

	if inRegistry(cmd) {
		fmt.Fprintln(w, cmd+" is a shell builtin")
		return
	}

	// Checking if a filename exists
	found := false
	isExec := false

	for _, dir := range paths {
		fullPath := filepath.Join(dir, words[1])
		fi, err := os.Stat(fullPath)
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
			fmt.Fprintf(w, "%s is %s\n", cmd, fullPath)
			return

		}
	}
	if !found {
		fmt.Fprintln(w, cmd+": not found")
	}
}
