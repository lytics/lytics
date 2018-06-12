package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "account",
		Usage:    "Account Info",
		Category: "Management API",
		Subcommands: []cli.Command{
			{
				Name:   "get",
				Usage:  "Show details of current authenticated account",
				Action: accountShow,
			},
			{
				Name:   "list",
				Usage:  "List accounts",
				Action: accountList,
			},
		},
	})
	addCommand(cli.Command{
		Name:     "accountuser",
		Usage:    "Account Admin User Info",
		Category: "Management API",
		Subcommands: []cli.Command{
			{
				Name:      "get",
				Usage:     "Show details of single account-user",
				ArgsUsage: "[id or email of user]",
				Action:    accountUserGet,
			},
			{
				Name:   "list",
				Usage:  "List account-users",
				Action: accountUserList,
			},
		},
	})
}
func accountShow(c *cli.Context) error {
	accountId := apikey
	if len(c.Args()) == 1 {
		accountId = c.Args().First()
	}
	item, err := client.GetAccount(accountId)
	exitIfErr(err, "Could not get account %q from api", accountId)
	resultWrite(c, &item)
	return nil
}
func accountList(c *cli.Context) error {
	items, err := client.GetAccounts()
	exitIfErr(err, "Could not get account list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list)
	return nil
}
func accountUserGet(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetUser(id)
	exitIfErr(err, "Could not get admin-user %s", id)
	resultWrite(c, &item)
	return nil
}
func accountUserList(c *cli.Context) error {
	items, err := client.GetUsers()
	exitIfErr(err, "Could not get account list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list)
	return nil
}
