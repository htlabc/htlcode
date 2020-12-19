package RvoteParam

type RvoteParam struct {
	ServerId string
	Term     int64
	//请求投票候选人id
	CandidateId string
	//候选人的最后日志条目索引值
	LastLogIndex int64
	//候选人最后日志条目任期
	LastLogTerm int64
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
	r.CandidateId = val
	return r
}

func (r *RvoteParam) SetCandidateId(val string) *RvoteParam {
	r.CandidateId = val
	return r
}

//候选人的最后日志条目索引值
func (r *RvoteParam) SetLastLogIndex(val int64) *RvoteParam {
	r.LastLogIndex = val
	return r

}

//候选人最后日志条目任期
func (r *RvoteParam) SetLastLogTerm(val int64) *RvoteParam {
	r.LastLogTerm = val
	return r
}
