package command

func (c *Cli) getWorks(id interface{}) (interface{}, error) {
	if id != nil && id != "" {
		work, err := c.Client.GetWork(id.(string), false)
		if err != nil {
			return nil, err
		}

		return work, nil
	} else {
		works, err := c.Client.GetWorks()
		if err != nil {
			return nil, err
		}

		return works, nil
	}
}
