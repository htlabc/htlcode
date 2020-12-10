package inter

import (
	entry "htl/myraft.com/entry"
)

type StateMachine interface {
	Apply(LogEntry *entry.LogEntry)
	Get(key string) *entry.LogEntry
	GetString(key string) string
	SetString(key string, value string)
	DelString(key ...string)
}
