package command

import (
	"fmt"
	"strings"
)

func schemaCommands(api *apiCommand) map[string]*command {
	c := &schema{apiCommand: api}
	return map[string]*command{
		"":       &command{c.HelpSchema, c.Schema, "Schema Show fields for a table."},
		"tables": &command{c.HelpTables, c.Schema, "Schema show tables."},
	}
}

type schema struct {
	*apiCommand
}

func (c *schema) HelpTables() string {
	helpText := fmt.Sprintf(`
Usage: lytics schema tables [options]

  List schema tables

%s
`, globalHelp)
	return strings.TrimSpace(helpText)
}
func (c *schema) HelpSchema() string {
	helpText := fmt.Sprintf(`
Usage: lytics schema [options]

  Get Schema and show columns

%s

Options:

`, globalHelp)
	return strings.TrimSpace(helpText)
}

func (c *schema) Schema(args []string) int {

	c.init(args, c.HelpSchema)
	table := c.f.Arg(0)
	if table == "" {
		table = "user"
	}
	//{As:"country", IsBy:false, Type:"string", ShortDesc:"Country", LongDesc:"Country Code", Froms:[]string{"default", "app", "clearbit_users"}, Identities:[]string{"geocountry", "GeoCountryCode"}}
	c.cols = []string{"as", "type", "shortdesc"}

	schema, err := c.l.GetSchemaTable(table)
	c.exitIfErr(err, "Could not get schema")

	items := make([]interface{}, len(schema.Columns))
	for i, col := range schema.Columns {
		items[i] = col
	}
	c.writeList(items)
	return 0
}
