package rpc

type RpcClient interface {
	Send(r Request) *Response
	SendTimeOut(r Request, timeout int) *Response
}
