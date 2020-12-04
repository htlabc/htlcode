package raft_client

const PUT = 0
const GET = 1

type ClientKVReq struct {
	reqType int
	key     string
	value   string
}

func (c *ClientKVReq) SetReqType(value int) {
	c.reqType = value
}

func (c *ClientKVReq) SetKey(value string) {
	c.key = value
}

func (c *ClientKVReq) SetValue(value string) {
	c.value = value
}

func (c *ClientKVReq) GetReqType() int {
	return c.reqType
}

func (c *ClientKVReq) GetKey() string {
	return c.key
}

func (c *ClientKVReq) GetValue() string {
	return c.value
}
