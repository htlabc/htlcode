package inter

import (
	"fmt"
	"htl/myraft.com/command"
	"htl/myraft.com/impl"
	"htl/myraft.com/memchange"
	"htl/myraft.com/rpc"
)

type ClusterMembershipChanges interface {
	AddPeer(newPeer *command.Peer) *memchange.ClusterMemberChageResult
	RemovePeer(oldPeer *command.Peer) *memchange.ClusterMemberChageResult
}

type ClusterMembershipChangesImpl struct {
	node impl.DefaultNode
}

func (c *ClusterMembershipChangesImpl) AddPeer(newPeer *command.Peer) *memchange.ClusterMemberChageResult {
	//检查新添加的peer是否已经存在集群内
	set := c.node.GetPeerSet()
	for x := set.GetPeers().Front(); x != nil; x = x.Next() {
		peer := x.Value.(command.Peer)
		if peer.GetAddr() == newPeer.GetAddr() {
			return nil
		}
	}
	if c.node.GetNodeStatus() == command.LEADER {
		c.node.GetMatchedIndexes()[*newPeer] = 0
		c.node.GetNextIndexes()[*newPeer] = 0
		for i := 0; i <= int(c.node.LogModule.GetLastIndex()); i++ {
			log := c.node.LogModule.Read(int64(i))
			c.node.Replication(*newPeer, log)
		}

		l := c.node.GetPeerSet().GetPeersWithOutSelf()

		for peer := l.Front(); peer.Value != nil; peer = peer.Next() {
			request := rpc.NewRequest().SetCmd(rpc.CHANGE_CONFIG_ADD).SetObj(newPeer).SetUrl(newPeer.GetAddr())
			response := c.node.RpcClient.Send(*request)
			result := response.GetResult().(*memchange.ClusterMemberChageResult)
			if result != nil && result.GetStatus() == memchange.SUCCESS {
				fmt.Printf("replication config success, peer : %s, newServer : %s", newPeer, newPeer.GetAddr())
			} else {
				fmt.Printf("replication config failed, peer : %s, newServer : %s ", newPeer, newPeer.GetAddr())
			}
		}
	}

	return nil
}

func (c *ClusterMembershipChangesImpl) RemovePeer(oldPeer *command.Peer) *memchange.ClusterMemberChageResult {
	//检查需要删除的peer是否存在集群内
	set := c.node.GetPeerSet()
	set.RemovePeer(oldPeer)
	if c.node.GetNodeStatus() == command.LEADER {
		delete(c.node.GetMatchedIndexes(), *oldPeer)
		delete(c.node.GetNextIndexes(), *oldPeer)
		l := c.node.GetPeerSet().GetPeersWithOutSelf()
		for peer := l.Front(); peer.Value != nil; peer = peer.Next() {
			request := rpc.NewRequest().SetCmd(rpc.CHANGE_CONFIG_REMOVE).SetObj(oldPeer).SetUrl(oldPeer.GetAddr())
			response := c.node.RpcClient.Send(*request)
			result := response.GetResult().(*memchange.ClusterMemberChageResult)
			if result != nil && result.GetStatus() == memchange.SUCCESS {
				fmt.Printf("remove config success, peer : %s, newServer : %s", oldPeer, oldPeer.GetAddr())
			} else {
				fmt.Printf("remove config failed, peer : %s, newServer : %s ", oldPeer, oldPeer.GetAddr())
			}
		}
	}
	return nil
}
