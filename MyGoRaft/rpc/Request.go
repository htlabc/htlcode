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

func (r *Request) SetCmd(value int) {
	r.CMD = value
}

func (r *Request) SetUrl(value string) {
	r.URL = value
}

func (r *Request) SetObj(value interface{}) {
	r.OBJ = value
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
