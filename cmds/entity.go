package cmds

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/araddon/gou"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "entity",
		Usage:    "Entity API:  Read a single User (or other table entity type) from a Table.",
		Category: "Data API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "table",
				Usage: "Table that describes the fields of this entity type/table.",
				Value: "user",
			},
		},
		Subcommands: []*cli.Command{
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
	if c.NArg() == 0 {
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
	exitIfErr(err, "eould not get entity %q from API", id)

	switch outputFormat {
	case "table":
		cols := make([]string, 0, len(ent.Fields))
		by := make(map[string]bool)
		for k := range ent.Fields {
			cols = append(cols, strings.ToLower(k))
		}
		byfields := gou.JsonHelper(ent.Meta).Strings("by_fields")
		for _, byfield := range byfields {
			by[byfield] = true
		}
		sort.Strings(cols)

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetHeader([]string{"field", "by", "value"})
		for _, col := range cols {
			byf := ""
			if _, isBy := by[col]; isBy {
				byf = "*"
			}

			if val, ok := ent.Fields[col]; ok {
				table.Append(rowToString([]interface{}{col, byf, val}))
			}
		}
		table.SetAutoFormatHeaders(false)
		table.Render()
		fmt.Println(tableString.String())

	case "json":
		jsonOut, err := json.MarshalIndent(ent.Fields, "", "  ")
		exitIfErr(err, "eould not marshal JSON")
		fmt.Printf("%s\n", string(jsonOut))
	default:
		resultWrite(c, &ent.Fields, "entity")
	}

	return nil
}
