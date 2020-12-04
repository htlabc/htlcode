package rpc

type RpcServer interface {
	Start(stopchan <-chan struct{}) error
	Stop(stopchan <-chan struct{}) error
	HandlerRequest(request Request, response *Response) error
}
