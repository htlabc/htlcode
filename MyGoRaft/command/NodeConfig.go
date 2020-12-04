package command

type NodeConfig struct {
	/** 自身 selfPort */
	selfPort int
	/** 所有节点地址. */
	peerAddrs []string
}
