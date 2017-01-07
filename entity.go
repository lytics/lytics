package main

func (c *Cli) getEntity(entitytype, fieldname, fieldval string, fields []string) (interface{}, error) {
	entity, err := c.Client.GetEntity(entitytype, fieldname, fieldval, fields)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
