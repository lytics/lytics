package main

import ()

func (c *Cli) getSchema(table string) (interface{}, error) {
	if table == "" {
		table = "user"
	}

	schema, err := c.Client.GetSchemaTable(table)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
