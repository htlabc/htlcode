package inter

import "htl/myraft.com/entry"

type LogModule interface {
	Write(logEntry *entry.LogEntry)
	Read(index int64) *entry.LogEntry
	RemoveOnstartIndex(startIndex int64)
	GetLast() *entry.LogEntry
	GetLastIndex() int64
}
