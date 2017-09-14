package command

import (
	"fmt"
	"strings"
)

func userCommands(api *apiCommand) map[string]*command {
	c := &user{apiCommand: api}
	return map[string]*command{
		"list": &command{c.HelpList, c.List, "User List, admi users for account."},
		"show": &command{c.HelpGet, c.Get, "Admin User Show Summary."},
	}
}

type user struct {
	*apiCommand
}

func (c *user) HelpList() string {
	helpText := fmt.Sprintf(`
Usage: lytics user list [options]

  List users

%s
`, globalHelp)
	return strings.TrimSpace(helpText)
}
func (c *user) HelpGet() string {
	helpText := fmt.Sprintf(`
Usage: lytics user show [options] id

  Get Admin User and show summary

%s

Options:

`, globalHelp)
	return strings.TrimSpace(helpText)
}

func (c *user) Get(args []string) int {

	c.init(args, c.HelpGet)
	id := c.f.Arg(0)
	if id == "" {
		c.ui.Error("Must provide user ID")
	}
	c.cols = []string{"email", "name", "created", "roles"}

	user, err := c.l.GetUser(id)
	c.exitIfErr(err, "Could not get user")
	// for _, ua := range  user.Accounts {
	// }

	c.writeSingle(user)
	return 0
}

func (c *user) List(args []string) int {
	c.init(args, c.HelpList)
	c.cols = []string{"email", "name", "created", "roles"}

	users, err := c.l.GetUsers()
	if err != nil {
		c.ui.Error(fmt.Sprintf("Could not get users %v", err))
		return 1
	}
	items := make([]interface{}, len(users))
	for i, u := range users {
		items[i] = u
	}
	c.writeList(items)
	return 0
}
