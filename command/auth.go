package command

import (
	"fmt"
	"strings"
)

func authCommands(api *apiCommand) map[string]*command {
	a := &auth{apiCommand: api}
	return map[string]*command{
		"list": &command{a.HelpList, a.List, "Auth List, 3rd party auth tokens."},
		"show": &command{a.HelpGet, a.Get, "Auth Show Summary."},
	}
}

type auth struct {
	*apiCommand
}

func (c *auth) HelpList() string {
	helpText := fmt.Sprintf(`
Usage: lytics auth list [options]

  List auths

%s
`, globalHelp)
	return strings.TrimSpace(helpText)
}
func (c *auth) HelpGet() string {
	helpText := fmt.Sprintf(`
Usage: lytics auth show [options] id

  Get Auth and show summary

%s

Options:

`, globalHelp)
	return strings.TrimSpace(helpText)
}

func (c *auth) Get(args []string) int {

	c.init(args, c.HelpGet)
	id := c.f.Arg(0)
	if id == "" {
		c.ui.Error("Must provide auth ID")
	}
	c.cols = []string{"id", "name", "created", "description"}

	auth, err := c.l.GetAuth(id)
	c.exitIfErr(err, "Could not get auth")

	c.writeSingle(auth)
	return 0
}

func (c *auth) List(args []string) int {
	c.init(args, c.HelpList)
	c.cols = []string{"id", "name", "created", "description"}

	auths, err := c.l.GetAuths()
	if err != nil {
		c.ui.Error(fmt.Sprintf("Could not get auths %v", err))
		return 1
	}
	items := make([]interface{}, len(auths))
	for i, auth := range auths {
		items[i] = auth
	}
	c.writeList(items)
	return 0
}
