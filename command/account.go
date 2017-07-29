package command

import (
	"fmt"
	"strings"
)

func accountCommands(api *apiCommand) map[string]*command {
	a := &account{apiCommand: api}
	return map[string]*command{
		"list": &command{a.HelpList, a.List, "Account List."},
		"show": &command{a.HelpGet, a.Get, "Account Show Summary."},
	}
}

type account struct {
	*apiCommand
}

func (c *account) HelpList() string {
	helpText := fmt.Sprintf(`
Usage: lytics account list [options]

  List accounts

%s

Options:

  -to=yesterday   End date to extract events.
  -format=json    Choose export format between json/csv.
  -event=E        Extract data for only event E.
  -out=STDOUT     Decides where to write the data.
`, globalHelp)
	return strings.TrimSpace(helpText)
}
func (c *account) HelpGet() string {
	helpText := fmt.Sprintf(`
Usage: lytics account show [options] id

  Get Account and show summary

%s

Options:

`, globalHelp)
	return strings.TrimSpace(helpText)
}

func (c *account) Get(args []string) int {
	//c.f.StringVar(&c.stuff, "stuff", "", "Account stuff")
	c.ui.Info("GET")
	c.init(args, c.HelpGet)
	id := c.f.Arg(0)
	if id == "" {
		c.ui.Error("Must provide account ID")
	}
	c.cols = []string{"id", "name", "created", "description"}

	//c.ui.Info(fmt.Sprintf("hello world aid=%v arg=%v", c.aid, c.f.Arg(0)))

	acct, err := c.l.GetAccount(id)
	c.exitIfErr(err, "Could not get account")

	c.writeSingle(acct)
	return 0
}

func (c *account) List(args []string) int {
	//c.f.StringVar(&c.stuff, "stuff", "", "Account stuff")
	c.init(args, c.HelpList)
	c.cols = []string{"id", "name", "created", "description"}

	c.ui.Info(fmt.Sprintf("Account List  id=%v", c.f.Arg(0)))

	accts, err := c.l.GetAccounts()
	if err != nil {
		c.ui.Error(fmt.Sprintf("Could not get accounts %v", err))
		return 1
	}
	items := make([]interface{}, len(accts))
	for i, acct := range accts {
		items[i] = acct
		//fmt.Printf("ACCT:  %T %#v \n", acct, acct)
	}
	c.writeList(items)
	return 0
}

// func items(accts []lytics.Account) {
// }
