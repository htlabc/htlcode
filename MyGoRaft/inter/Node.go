package inter

import (
	"htl/myraft.com/entry"
	vote "htl/myraft.com/entry/RvoteParam"
	raft_client "htl/myraft.com/raft.client"
)

type Node interface {
	handlerRequestVote(vote vote.RvoteParam) *entry.RvoteResult
	handlerAppendEntries(param entry.AppendEntryParam) *entry.AppendEntryResult
	handlerClientRequest(request raft_client.ClientKVReq) *raft_client.ClientKVAck
	redirect(request raft_client.ClientKVReq) *raft_client.ClientKVAck
}
