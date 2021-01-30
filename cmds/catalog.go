package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "schema",
		Usage:    "Schema (Catalog) information about Lytics Tables & Queries",
		Category: "Data API",
		Subcommands: []*cli.Command{
			{
				Name:  "tables",
				Usage: "API of tables that make up schema",
				Subcommands: []*cli.Command{
					{
						Name:      "get",
						Usage:     "Show details of current requested table schema",
						UsageText: "Get Detail of Single Table Schema",
						ArgsUsage: "[table name]",
						Action:    schemaTableGet,
					},
					{
						Name:   "list",
						Usage:  "List Tables",
						Action: schemaTableList,
					},
				},
			},
			{
				Name:  "queries",
				Usage: "API of queries that make up schema",
				Subcommands: []*cli.Command{
					{
						Name:      "get",
						Usage:     "Show details of single query by alias",
						UsageText: "Get Detail of Single Query",
						ArgsUsage: "[query alias aka slug]",
						Action:    schemaQueryGet,
					},
					{
						Name:   "list",
						Usage:  "List Queries",
						Action: schemaQueryList,
					},
					{
						Name:   "watch",
						Usage:  "Watch Queries",
						Action: schemaQueryWatch,
					},
				},
			},
		},
	})
}

func schemaTableGet(c *cli.Context) error {
	id := "user"
	if c.NArg() > 0 {
		id = c.Args().First()
	}
	item, err := client.GetSchemaTable(id)
	exitIfErr(err, "could not get schema %q from API", id)
	list := make([]lytics.TableWriter, len(item.Columns))
	for i, item := range item.Columns {
		val := item
		list[i] = &val
	}
	resultWrite(c, list, fmt.Sprintf("schema_table_%s", id))
	return nil
}
func schemaTableList(c *cli.Context) error {
	items, err := client.GetSchema()
	exitIfErr(err, "could not get schema tables list")
	list := make([]lytics.TableWriter, 0, len(items))
	for _, item := range items {
		list = append(list, item)
	}
	resultWrite(c, list, "schema_table_list")
	return nil
}
func schemaQueryGet(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetQuery(id)
	exitIfErr(err, "could not get schema query %q from API", id)
	resultWrite(c, &item, fmt.Sprintf("schema_query_%s", item.Id))
	return nil
}
func schemaQueryList(c *cli.Context) error {
	items, err := client.GetQueries()
	exitIfErr(err, "could not get schema queries list")
	list := make([]lytics.TableWriter, 0, len(items))
	for _, item := range items {
		list = append(list, item)
	}
	resultWrite(c, list, "schema_query_list")
	return nil
}
