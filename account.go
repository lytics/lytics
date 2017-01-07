package main

func (c *Cli) getAccounts(id interface{}) (interface{}, error) {
	if id != nil && id != "" {
		acct, err := c.Client.GetAccount(id.(string))
		if err != nil {
			return nil, err
		}

		return acct, nil
	} else {
		accts, err := c.Client.GetAccounts()
		if err != nil {
			return nil, err
		}

		return accts, nil
	}
}
