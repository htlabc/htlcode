package command

type NodeConfig struct {
	/** 自身 selfPort */
	SelfPort int
	/** 所有节点地址. */
	PeerAddrs []string
}
