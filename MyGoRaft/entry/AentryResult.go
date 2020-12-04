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

func (a *AppendEntryResult) aentryResult(success bool) {
	a.success = success
}

func (a *AppendEntryResult) AentryResult(term int64, success bool) {
	a.success = success
	a.term = term
}

func (a *AppendEntryResult) fail() *AppendEntryResult {
	a.success = false
	return a
}

func (a *AppendEntryResult) ok() *AppendEntryResult {
	a.success = true
	return a
}
