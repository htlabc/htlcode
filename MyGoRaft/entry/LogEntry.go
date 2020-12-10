package entry

type LogEntry struct {
	index   int64
	term    int64
	command *Command
}

func (logentry *LogEntry) Error() string {
	panic("implement me")
}

func (entry *LogEntry) GetCommand() *Command {
	return entry.command
}

func (entry *LogEntry) GetTerm() int64 {
	return entry.term
}

func (logentry *LogEntry) LogEntry(entryBuilder LogEntryBuilder) *LogEntry {
	return &LogEntry{index: entryBuilder.term, term: entryBuilder.term, command: entryBuilder.command}
}

func (entry *LogEntry) GetIndex() int64 {
	return entry.index
}

func (entry *LogEntry) SetIndex(val int64) {
	entry.index = val
}

type LogEntryBuilder struct {
	index   int64
	term    int64
	command *Command
}

func (l *LogEntryBuilder) Index(val int64) *LogEntryBuilder {
	l.index = val
	return l
}

func (l *LogEntryBuilder) Term(val int64) *LogEntryBuilder {
	l.term = val
	return l
}

func (l *LogEntryBuilder) Command(c *Command) *LogEntryBuilder {
	l.command = c
	return l
}

func (l *LogEntryBuilder) LogEntryBuild() *LogEntry {
	log := new(LogEntry)
	log.LogEntry(*l)
	return log
}
