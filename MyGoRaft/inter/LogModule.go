package inter

import "htl/myraft.com/entry"

type LogModule interface {
	write(Logentry entry.LogEntry)
	Read(index int64)
	removeOnstartIndex(startIndex int64)
	GetLast() *entry.LogEntry
	GetLastIndex() int64
}
