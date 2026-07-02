package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Running a program
func handleProgram(w io.Writer, cmd string, args []string) {
	program := exec.Command(cmd, args...)
	program.Stdout = w
	program.Stdin = os.Stdin
	program.Stderr = os.Stderr

	// Executing the child process
	err := program.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
