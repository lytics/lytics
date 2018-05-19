package cmds

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "account",
		Usage:    "Account Info",
		Category: "Management API",
		Action: func(c *cli.Context) error {
			fmt.Println("account: ", c.Args().First())
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:   "show",
				Usage:  "Show details of current authenticated account",
				Action: accountShow,
			},
			{
				Name:  "list",
				Usage: "List accounts",
				Action: func(c *cli.Context) error {
					fmt.Println("list account: ", c.Args().First())
					return nil
				},
			},
		},
	})
}
func accountShow(c *cli.Context) error {
	fmt.Println("show account ", c.Args().First())
	acct, err := client.GetAccount(c.Args().First())
	if err != nil {
		fmt.Println("error getting account", err)
		return err
	}
	fmt.Println("%+v", acct)
	return nil
}

func (c *Cli) getAccounts(id interface{}) (interface{}, error) {
	accts, err := c.Client.GetAccounts()
	if err != nil {
		return nil, err
	}

	return accts, nil
}
