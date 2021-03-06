package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(cli.Command{
		Name:     "account",
		Usage:    "Account Info",
		Category: "Management API",
		Subcommands: []*cli.Command{
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
		Subcommands: []*cli.Command{
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
	accountID := apikey
	if c.NArg() == 1 {
		accountID = c.Args().First()
	}
	item, err := client.GetAccount(accountID)
	exitIfErr(err, "could not get account %q from API", accountID)
	resultWrite(c, &item, fmt.Sprintf("account_%s", item.Name))
	return nil
}
func accountList(c *cli.Context) error {
	items, err := client.GetAccounts()
	exitIfErr(err, "could not get account list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list, "account_list")
	return nil
}
func accountUserGet(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetUser(id)
	exitIfErr(err, "could not get admin-user %s", id)
	resultWrite(c, &item, fmt.Sprintf("account_%s", item.Name))
	return nil
}
func accountUserList(c *cli.Context) error {
	items, err := client.GetUsers()
	exitIfErr(err, "could not get account list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list, "account_user_list")
	return nil
}
