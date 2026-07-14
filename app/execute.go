package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func executeCmd(words Command) error {
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr
	var files []*os.File

	for _, redir := range words.Redirections {
		flags := os.O_RDWR | os.O_CREATE
		if redir.Operator == ">>" {
			flags |= os.O_APPEND
		} else if redir.Operator == ">" {
			flags |= os.O_TRUNC
		}

		// Creating the parent directories
		dir := filepath.Dir(redir.Filename)
		if dir != "." && dir != "/" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
		// Configuring the child process
		file, err := os.OpenFile(redir.Filename, flags, 0666)
		if err != nil {
			return err
		}
		// Tracking the files we have opened to close later
		files = append(files, file)

		switch redir.FD {
		case 1:
			// Writes the child process directly to the file redirected
			stdout = file
		case 2:
			stderr = file
		}
	}
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	cmdName := words.Name
	args := words.Args

	input := append([]string{cmdName}, args...)

	// Builtin commands
	handler, ok := builtinRegistry[cmdName]
	if ok {
		handler(stdout, input)
		return nil
	}
	// External commands
	program := exec.Command(cmdName, args...)
	program.Stdout = stdout
	program.Stdin = os.Stdin
	program.Stderr = stderr

	// Simply return the error - let the caller decide how to handle it
	return program.Run()
}
