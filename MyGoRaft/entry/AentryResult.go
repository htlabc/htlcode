package entry

type AppendEntryResult struct {
	term    int64
	success bool
}

func (a *AppendEntryResult) GetTerm() int64 {
	return a.term
}

func (a *AppendEntryResult) GetStatus() bool {
	return a.success
}

func (a *AppendEntryResult) Builder(term int64, success bool) {
	a.success = success
	a.term = term
}

func (a *AppendEntryResult) Fail() *AppendEntryResult {
	a.success = false
	return a
}

func (a *AppendEntryResult) Ok() *AppendEntryResult {
	a.success = true
	return a
}
