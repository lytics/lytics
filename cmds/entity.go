package cmds

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/apcera/termtables"
	"github.com/araddon/gou"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "entity",
		Usage:    "Entity Api:  Read a single User (or other table entity type) from a Table.",
		Category: "Data API",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "table",
				Usage: "table that describes the fields of this entity type/table.",
				Value: "user",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:      "get",
				Usage:     "Show full entity detail",
				ArgsUsage: "[id of entity in format:   fieldname/fieldval]",
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
	if table == "" {
		table = "user"
	}
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return fmt.Errorf("expected:    fieldname/value    with / slash in between got %q", id)
	}

	ent, err := client.GetEntity(table, parts[0], parts[1], nil)
	exitIfErr(err, "Could not get entity %q from api", id)

	switch outputFormat {
	case "table":
		cols := make([]string, 0, len(ent.Fields))
		by := make(map[string]bool)
		for k, _ := range ent.Fields {
			cols = append(cols, strings.ToLower(k))
		}
		byfields := gou.JsonHelper(ent.Meta).Strings("by_fields")
		for _, byfield := range byfields {
			by[byfield] = true
		}
		sort.Strings(cols)

		table := termtables.CreateTable()
		headers := []interface{}{"field", "by", "value"}
		table.AddHeaders(headers...)
		for _, col := range cols {
			byf := ""
			if _, isBy := by[col]; isBy {
				byf = "*"
			}

			if val, ok := ent.Fields[col]; ok {
				table.AddRow(col, byf, val)
			}
		}
		fmt.Println(table.Render())

	case "json":
		jsonOut, err := json.MarshalIndent(ent.Fields, "", "  ")
		exitIfErr(err, "Could not marshal json")
		fmt.Printf("%s\n", string(jsonOut))
	default:
		resultWrite(c, &ent.Fields)
	}

	return nil
}
