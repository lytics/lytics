package cmds

func (c *Cli) getQueries(alias string) (interface{}, error) {
	if alias == "" {
		return c.Client.GetQueries()
	}
	return c.Client.GetQuery(alias)
}
