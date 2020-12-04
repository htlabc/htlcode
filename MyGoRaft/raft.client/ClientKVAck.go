package raft_client

type ClientKVAck struct {
	object interface{}
}

func (c *ClientKVAck) ok() *ClientKVAck {
	c.object = "ok"
	return c
}

func (c *ClientKVAck) fail() *ClientKVAck {
	c.object = "fail"
	return c
}
