package impl

import (
	"htl/myraft.com/entry"
	"htl/myraft.com/entry/RvoteParam"
	inter "htl/myraft.com/inter"
	"sync"
)

type DefaultConsensus struct {
	node       DefaultNode
	voteRock   sync.Mutex
	appendLock sync.Mutex
	inter.Consenus
}

func NewDefaultConsenus() *DefaultConsensus {
	return &DefaultConsensus{}
}

func RequestVote(param RvoteParam.RvoteParam) *entry.RvoteResult {
	return nil
}

func AppendEntries(param entry.AppendEntryParam) *entry.AppendEntryResult {
	return nil
}
