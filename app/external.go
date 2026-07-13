package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Running a program
func handleProgram(stdout io.Writer, stderr io.Writer, cmd string, args []string) {
	program := exec.Command(cmd, args...)
	program.Stdout = stdout
	program.Stdin = os.Stdin
	program.Stderr = stderr

	// Executing the child process
	err := program.Run()
	if err != nil {
		// Checking the type error
		var exitError *exec.ExitError
		var execError *exec.Error
		if errors.As(err, &execError) {
			fmt.Fprintf(os.Stderr, "%s: not found\n", cmd)
		} else if errors.As(err, &exitError) {
			return
		} else {
			fmt.Fprintf(os.Stderr, "%s: unexpected error: %v\n", cmd, err)
		}

	}
}
