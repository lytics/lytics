package command

func (c *Cli) getUsers(id interface{}) (interface{}, error) {
	if id != nil && id != "" {
		user, err := c.Client.GetUser(id.(string))
		if err != nil {
			return nil, err
		}

		return user, nil
	} else {
		users, err := c.Client.GetUsers()
		if err != nil {
			return nil, err
		}

		return users, nil
	}
}
