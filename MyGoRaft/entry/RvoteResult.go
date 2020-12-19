package entry

type RvoteResult struct {
	term        int64
	voteGranted bool
}

func NewRvoteResult() *RvoteResult {
	return &RvoteResult{}
}

func (r *RvoteResult) SetTerm(val int64) *RvoteResult {
	r.term = val
	return r
}

func (r *RvoteResult) SetVoteGranted(val bool) *RvoteResult {
	r.voteGranted = val
	return r
}

func (r *RvoteResult) Term() int64 {
	return r.term
}

func (r *RvoteResult) VoteGranted() bool {
	return r.voteGranted
}
