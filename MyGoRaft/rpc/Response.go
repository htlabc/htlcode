package rpc

type Response struct {
	Result interface{}
}

func NewResponse(value interface{}) *Response {
	r := new(Response)
	r.Result = value
	return r
}

func (r *Response) SetResult(value interface{}) {
	r.Result = value
}

func (r *Response) GetResult() interface{} {
	return r.Result
}
