package main

import (
	"flag"
	"htl/myraft.com/command"
	"htl/myraft.com/impl"
	"sync"
)

func main() {

	serverport := flag.Int("SERVERPORT", 8001, "raft server port")
	flag.Parse()
	var peerAddr []string = []string{"localhost:8775", "localhost:8776", "localhost:8777", "localhost:8778", "localhost:8779"}
	nodeconfig := command.NodeConfig{}
	nodeconfig.SelfPort = *serverport
	nodeconfig.PeerAddrs = peerAddr
	node := impl.GetInstance()
	node.SetNodeConfig(nodeconfig)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		node.Init()
		wg.Done()
	}()

	wg.Done()

}
