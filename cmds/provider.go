package cmds

func (c *Cli) getProviders(id interface{}) (interface{}, error) {
	if id != nil && id != "" {
		provider, err := c.Client.GetProvider(id.(string))
		if err != nil {
			return nil, err
		}

		return provider, nil
	} else {
		providers, err := c.Client.GetProviders()
		if err != nil {
			return nil, err
		}

		return providers, nil
	}
}
