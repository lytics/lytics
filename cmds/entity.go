package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "entity",
		Usage:    "Entity Api",
		Category: "Data API",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "table",
				Usage: "table that describes this entity type",
				Value: "user",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:      "get",
				Usage:     "Show details of entity",
				ArgsUsage: "[id of entity in format   fieldname/fieldval]",
				Action:    entityGet,
			},
		},
	})
}
func entityGet(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return fmt.Errorf("expected one arg (field/fieldval)")
	}
	id := c.Args().First()
	table := c.String("table")
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return fmt.Errorf("expected    fieldname/value    with / slash in between got %q", id)
	}
	item, err := client.GetEntity(table, parts[0], parts[2], nil)
	exitIfErr(err, "Could not get auth %q from api", id)
	resultWrite(c, &item)
	return nil
}
