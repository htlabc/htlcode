package inter

import (
	"htl/myraft.com/command"
	"htl/myraft.com/memchange"
)

type ClusterMembershipChanges interface {
	AddPeer(newPeer command.Peer) memchange.ClusterMemberChageResult
	RemovePeer(oldPeer command.Peer) memchange.ClusterMemberChageResult
}
