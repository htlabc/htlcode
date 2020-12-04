package inter

import (
	entry "htl/myraft.com/entry"
)

type StateMachine interface {
	Apply(bucket string, LogEntry entry.LogEntry)
	Get(bucket string, key string) entry.LogEntry
	GetString(bucket string, key string) string
	SetString(bucket string, key string, value string)
	DelString(key ...string)
}
