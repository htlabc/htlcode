package entry

type RvoteResult struct {
	term        uint64
	voteGranted bool
}

func (r *RvoteResult) SetTerm(val uint64) {
	r.term = val
}

func (r *RvoteResult) SetVoteGranted(val bool) {
	r.voteGranted = val
}

func (r *RvoteResult) Term() uint64 {
	return r.term
}

func (r *RvoteResult) VoteGranted() bool {
	return r.voteGranted
}
