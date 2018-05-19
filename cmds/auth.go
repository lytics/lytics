package cmds

func (c *Cli) getAuths(id interface{}) (interface{}, error) {
	if id != nil && id != "" {
		auth, err := c.Client.GetAuth(id.(string))
		if err != nil {
			return nil, err
		}

		return auth, nil
	} else {
		auths, err := c.Client.GetAuths()
		if err != nil {
			return nil, err
		}

		return auths, nil
	}
}
