package terminal

import (
	//"log"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

// Creating the error types
var (
	ErrInterrupt = &TerminalError{"interrupt"}
	ErrEOF       = &TerminalError{"EOF"}
)

type TerminalError struct {
	msg string
}

func (e *TerminalError) Error() string {
	return e.msg
}

// Configuring tab completer choices
var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
)

type Terminal struct {
	rl *readline.Instance
}

func New() (*Terminal, error) {
	read, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}
	return &Terminal{rl: read}, err
}

func (t *Terminal) Close() error {
	return t.rl.Close()
}

func (t *Terminal) ReadCommand() (string, error) {
	line, err := t.rl.Readline()
	if err != nil {
		if err == readline.ErrInterrupt {
			return "", ErrInterrupt
		}
		if err == io.EOF {
			return "", ErrEOF
		}
		return "", err
	}
	// Removes the newline at the end
	return strings.TrimSpace(line), err
}
