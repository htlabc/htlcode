package rpc

const R_VOTE = 0
const A_ENTRIS = 1
const CLIENT_REQ = 2
const CHANGE_CONFIG_ADD = 3
const CHANGE_CONFIG_REMOVE = 4

type Request struct {
	CMD int
	OBJ interface{}
	URL string
}

func (r *Request) SetCmd(value int) *Request {
	r.CMD = value
	return r
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) SetUrl(value string) *Request {
	r.URL = value
	return r

}

func (r *Request) SetObj(value interface{}) *Request {
	r.OBJ = value
	return r

}

func (r *Request) GetObj() interface{} {
	return r.OBJ
}

func (r *Request) GetUrl() string {
	return r.URL

}

func (r *Request) GetCmd() int {
	return r.CMD

}
