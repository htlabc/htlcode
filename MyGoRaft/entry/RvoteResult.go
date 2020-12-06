package entry

type RvoteResult struct {
	term        int64
	voteGranted bool
}

func (r *RvoteResult) SetTerm(val int64) {
	r.term = val
}

func (r *RvoteResult) SetVoteGranted(val bool) {
	r.voteGranted = val
}

func (r *RvoteResult) Term() int64 {
	return r.term
}

func (r *RvoteResult) VoteGranted() bool {
	return r.voteGranted
}
