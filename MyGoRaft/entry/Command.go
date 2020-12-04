package entry

type Command struct {
	key   string
	value string
}

func (c *Command) SetKey(key string) {
	c.key = key
}

func (c *Command) SetValue(value string) {
	c.value = value
}

func (c *Command) GetKey() string {
	return c.key
}

func (c *Command) GetValue() string {
	return c.value
}

func Build(key string, value string) *Command {
	return &Command{key: key, value: value}
}
