package raft_client

type ClientKVAck struct {
	object interface{}
}

func NewClientKVAck() *ClientKVAck {
	return &ClientKVAck{}
}

func (c *ClientKVAck) Ojbect(obj interface{}) *ClientKVAck {
	c.object = obj
	return c
}

func (c *ClientKVAck) Ok() *ClientKVAck {
	c.object = "ok"
	return c
}

func (c *ClientKVAck) Fail() *ClientKVAck {
	c.object = "fail"
	return c
}
