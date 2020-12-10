package entry

type Command struct {
	key   string
	value string
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) SetKey(key string) *Command {
	c.key = key
	return c
}

func (c *Command) SetValue(value string) *Command {
	c.value = value
	return c
}

func (c *Command) GetKey() string {
	return c.key
}

func (c *Command) GetValue() string {
	return c.value
}
