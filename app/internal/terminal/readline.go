package terminal

import (
	//"log"
	//"io"
	"strings"

	"github.com/chzyer/readline"
)

/**
// Configuring tab completer choices
var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
)

// Replaces the bufio.NewScanner
// Initializing the readline instance
var l, err = readline.NewEx(&readline.Config{
	Prompt:          "$ ",
	HistoryFile:     "/tmp/readline.tmp",
	AutoComplete:    completer,
	InterruptPrompt: "^C",
	EOFPrompt:       "exit",
})

if err != nil {
	log.Fatalf("Error initializing readline: %v", err)
}
defer l.Close()
**/

func ReadCommand(l *readline.Instance) (string, error) {
	line, err := l.Readline()
	if err != nil {
		return "", err
	}
	// Removes the newline at the end
	return strings.TrimSpace(line), err
}
