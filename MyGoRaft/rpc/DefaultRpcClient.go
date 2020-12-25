package rpc

import (
	"fmt"
	"log"
	"net/rpc"
)

type DefaultRpcClient struct {
}

func NewDefaultRpcClient() *DefaultRpcClient {
	return &DefaultRpcClient{}
}

func (d *DefaultRpcClient) SendTimeOut(r Request, timeout int) *Response {
	return nil
}

func (d *DefaultRpcClient) Send(r Request) *Response {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8096")
	if err != nil {
		log.Fatalln()
	}
	var res Response
	err = conn.Call("DefaultRPCServer.HandleRaftRequest", r, &res)
	if err != nil {
		fmt.Println(err)
	}

	return &res

}
