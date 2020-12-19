package impl

import (
	"fmt"
	"htl/myraft.com/command"
	"htl/myraft.com/entry"
	"htl/myraft.com/entry/RvoteParam"
	inter "htl/myraft.com/inter"
	"sync"
	"time"
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

func (consensus DefaultConsensus) RequestVote(param *RvoteParam.RvoteParam) *entry.RvoteResult {

	revoteResult := entry.NewRvoteResult()
	consensus.voteRock.Lock()
	defer consensus.voteRock.Unlock()

	failedResult := revoteResult.SetTerm(consensus.node.currentTerm).SetVoteGranted(false)
	if param.Term < consensus.node.currentTerm {
		return failedResult
	}

	// (当前节点并没有投票 或者 已经投票过了且是对方节点) && 对方日志和自己一样新

	if consensus.node.voteFor == "" || consensus.node.voteFor == param.CandidateId {
		if consensus.node.LogModule.GetLast() != nil {
			if consensus.node.LogModule.GetLast().GetTerm() > param.Term {
				return failedResult
			}

			if consensus.node.LogModule.GetLastIndex() > param.LastLogIndex {
				return failedResult
			}
		}

		consensus.node.status = command.FOLLOWER
		peer := &command.Peer{}
		peer.SetAddr(param.CandidateId)
		consensus.node.peerset.SetLeader(peer)
		consensus.node.voteFor = param.ServerId
		return revoteResult.SetTerm(consensus.node.currentTerm).SetVoteGranted(true)
	}
	return nil
}

func (consenus *DefaultConsensus) AppendEntries(param *entry.AppendEntryParam) *entry.AppendEntryResult {
	result := &entry.AppendEntryResult{}
	consenus.appendLock.Lock()
	defer consenus.appendLock.Unlock()
	result.Builder(consenus.node.currentTerm, false)

	if param.GetTerm() < consenus.node.currentTerm {
		return result
	}

	consenus.node.preElectionTime = time.Now()
	consenus.node.preHeartBeatTime = time.Now()

	peer := &command.Peer{}
	peer.SetAddr(param.GetLeaderId())
	consenus.node.peerset.SetLeader(peer)

	if param.GetTerm() >= consenus.node.currentTerm {
		consenus.node.status = command.FOLLOWER
	}

	consenus.node.currentTerm = param.GetTerm()
	//心跳的entries为nil
	if param.GetEntries() == nil || len(param.GetEntries()) == 0 {
		fmt.Printf("node %s append heartbeat success , he's term : %s, my term : %s",
			param.GetLeaderId(), param.GetTerm(), consenus.node.currentTerm)
		heartbeatResult := &entry.AppendEntryResult{}
		heartbeatResult.Builder(consenus.node.currentTerm, true)
		return heartbeatResult
	}

	if consenus.node.LogModule.GetLastIndex() != 0 && param.GetPreLogIndex() != 0 {
		logentry := &entry.LogEntry{}
		logentry = consenus.node.LogModule.Read(param.GetPreLogIndex())
		if logentry != nil {
			if logentry.GetTerm() != param.GetTerm() {
				return result
			}
		}
	}

	// 如果已经存在的日志条目和新的产生冲突（索引值相同但是任期号不同），删除这一条和之后所有的
	existLog := &entry.LogEntry{}
	existLog = consenus.node.LogModule.Read(param.GetPreLogIndex() + 1)
	if existLog != nil && existLog.GetTerm() != param.GetEntries()[0].GetTerm() {
		consenus.node.LogModule.RemoveOnstartIndex(param.GetPreLogIndex() + 1)
	} else {
		return result.Ok()
	}

	//写日志进状态机

	for _, log := range param.GetEntries() {
		consenus.node.LogModule.Write(&log)
		consenus.node.stateMachine.Apply(&log)
		result.Ok()
	}

	//如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 新日志条目索引值中较小的一个

	if param.GetLeaderCommit() > consenus.node.commitIndex {
		if param.GetLeaderCommit() > consenus.node.LogModule.GetLastIndex() {
			consenus.node.commitIndex = consenus.node.LogModule.GetLastIndex()
			consenus.node.lastApplied = consenus.node.LogModule.GetLastIndex()
		} else {
			consenus.node.commitIndex = param.GetLeaderCommit()
			consenus.node.lastApplied = param.GetLeaderCommit()
		}

	}

	result.Builder(consenus.node.currentTerm, true)

	consenus.node.status = command.FOLLOWER

	return result

	return nil
}
