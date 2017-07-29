package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/lytics/lytics/command"
)

func main() {
	ui := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorBlue,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorGreen,
		Ui:          &cli.BasicUi{Writer: os.Stdout},
	}
	c := cli.NewCLI("lytics", "1.0.0")
	c.HelpFunc = helpFunc
	c.Args = os.Args[1:]
	c.Commands = command.Commands(ui)

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

// BasicHelpFunc generates some basic help output that is usually good enough
// for most CLI applications.
func helpFunc(commands map[string]cli.CommandFactory) string {
	var buf bytes.Buffer
	buf.WriteString("Usage: lytics [--version] [--help] <command> [<args>]\n\n")
	buf.WriteString("Available commands are:\n")

	// Get the list of keys so we can sort them, and also get the maximum
	// key length so they can be aligned properly.
	keys := make([]string, 0, len(commands))
	maxKeyLen := 0
	for key := range commands {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		keys = append(keys, key)
	}
	sort.Strings(keys)
	lastSub := ""

	for _, key := range keys {
		parts := strings.Split(key, " ")
		if len(parts) == 2 {
			if lastSub != parts[0] {
				buf.WriteString(fmt.Sprintf("\n%s\n", parts[0]))
			}
		} else {
			buf.WriteString(fmt.Sprintf("\n%s\n", key))
		}
		lastSub = parts[0]
		commandFunc, ok := commands[key]
		if !ok {
			// This should never happen since we JUST built the list of
			// keys.
			panic("command not found: " + key)
		}

		command, err := commandFunc()
		if err != nil {
			log.Printf("[ERR] cli: Command '%s' failed to load: %s",
				key, err)
			continue
		}

		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, command.Synopsis()))
	}

	return buf.String()
}
