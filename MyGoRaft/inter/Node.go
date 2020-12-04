package inter

import "../entry"
import "../entry/RvoteParam"
import "../raft.client"

type Node interface {
	handlerRequestVote(param RvoteParam.RvoteParam) entry.RvoteResult
	handlerAppendEntries(param entry.AentryParam) entry.AentryResult
	handlerClientRequest(request raft_client.ClientKVReq) raft_client.ClientKVAck
	redirect(request raft_client.ClientKVReq) raft_client.ClientKVAck
}
