package RvoteParam

type RvoteParam struct {
	ServerId string
	Term     int64
	//请求投票候选人id
	candidateId string
	//候选人的最后日志条目索引值
	lastLogIndex int64
	//候选人最后日志条目任期
	lastLogTerm int64
}

func NewRvoteParam() *RvoteParam {
	return &RvoteParam{}
}

func (r *RvoteParam) SetTerm(val int64) *RvoteParam {
	r.Term = val
	return r
}

func (r *RvoteParam) SetServerID(serverid string) *RvoteParam {
	r.ServerId = serverid
	return r
}

func (r *RvoteParam) SetCnadiateID(val string) *RvoteParam {
	r.candidateId = val
	return r
}

func (r *RvoteParam) SetCandidateId(val string) *RvoteParam {
	r.candidateId = val
	return r
}

//候选人的最后日志条目索引值
func (r *RvoteParam) SetLastLogIndex(val int64) *RvoteParam {
	r.lastLogIndex = val
	return r

}

//候选人最后日志条目任期
func (r *RvoteParam) SetLastLogTerm(val int64) *RvoteParam {
	r.lastLogTerm = val
	return r
}
