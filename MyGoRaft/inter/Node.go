package inter

import (
	"htl/myraft.com/entry"
	vote "htl/myraft.com/entry/RvoteParam"
	raft_client "htl/myraft.com/raft.client"
)

type Node interface {
	HandlerRequestVote(vote vote.RvoteParam) *entry.RvoteResult
	HandlerAppendEntries(param entry.AppendEntryParam) *entry.AppendEntryResult
	HandlerClientRequest(request raft_client.ClientKVReq) *raft_client.ClientKVAck
	Redirect(request raft_client.ClientKVReq) *raft_client.ClientKVAck
}
