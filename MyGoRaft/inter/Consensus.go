package inter

import (
	"htl/myraft.com/entry"
	"htl/myraft.com/entry/RvoteParam"
)

type Consenus interface {
	RequestVote(param RvoteParam.RvoteParam) entry.RvoteResult
	AppendEntries(param entry.AppendEntryParam) entry.AppendEntryResult
}
