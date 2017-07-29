package main

import (
	"log"
	"os"

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
	c.Args = os.Args[1:]
	c.Commands = command.Commands(ui)

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
