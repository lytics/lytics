package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "auth",
		Usage:    "Auth Token/Keys provided to Lytics",
		Category: "Management API",
		Subcommands: []cli.Command{
			{
				Name:        "get",
				Usage:       "Show details of current requested id auth (but not encrypted values)",
				UsageText:   "Get Detail of Single Auth",
				Description: "no really, there is a lot of details",
				ArgsUsage:   "[id of auth]",
				Action:      authGet,
			},
			{
				Name:   "list",
				Usage:  "List auths",
				Action: authList,
			},
		},
	})
}
func authGet(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetAuth(id)
	exitIfErr(err, "Could not get auth %q from api", id)
	resultWrite(c, &item)
	return nil
}
func authList(c *cli.Context) error {
	items, err := client.GetAuths()
	exitIfErr(err, "Could not get auths list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list)
	return nil
}
