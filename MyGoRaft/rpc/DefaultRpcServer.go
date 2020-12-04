package rpc

import (
	"fmt"
	"htl/myraft.com/command"
	"htl/myraft.com/entry"
	vote "htl/myraft.com/entry/RvoteParam"
	"htl/myraft.com/impl"
	raft_client "htl/myraft.com/raft.client"
	"net"
	"net/rpc"
)

var FLAG bool = false
var rpcserver *DefaultRpcServer

func NewDefaultRpcServer() *DefaultRpcServer {
	if FLAG == true {
		return rpcserver
	}
	rpcserver = &DefaultRpcServer{}
	return rpcserver
}

type DefaultRpcServer struct {
	Node *impl.DefaultNode
}

func (d *DefaultRpcServer) Start(stopchan <-chan struct{}) error {
	go func() {
		defaultrpcserver := new(DefaultRpcServer)
		rpc.Register(defaultrpcserver)
		listener, err := net.Listen("tcp", "127.0.0.1:8001")
		defer func() {
			listener.Close()
			<-stopchan
		}()
		if err != nil {
			fmt.Println(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			rpc.ServeConn(conn)
		}
	}()

	return nil

}

func (d *DefaultRpcServer) Stop(stopchan <-chan struct{}) error {
	return nil

}

func (rpc *DefaultRpcServer) HandlerRequest(request Request, response *Response) error {

	if request.GetCmd() == R_VOTE {
		r := request.GetObj().(vote.RvoteParam)
		response.SetResult(rpc.Node.HandlerRequestVote(r))
		return nil
	} else if request.GetCmd() == A_ENTRIS {
		r := request.GetObj().(entry.AppendEntryParam)
		response.SetResult(rpc.Node.HandlerAppendEntries(r))
		return nil
	} else if request.GetCmd() == CLIENT_REQ {
		r := request.GetObj().(raft_client.ClientKVReq)
		response.SetResult(rpc.Node.HandlerClientRequest(r))
		return nil
	} else if request.GetCmd() == CHANGE_CONFIG_REMOVE {
		req := request.GetObj().(command.Peer)
		response.SetResult(rpc.Node.AddPeers(req))
		return nil
	} else if request.GetCmd() == CHANGE_CONFIG_ADD {
		req := request.GetObj().(command.Peer)
		response.SetResult(rpc.Node.RemovePeers(req))
		return nil
	}
	return nil
}
