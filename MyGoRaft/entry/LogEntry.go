package entry

type LogEntry struct {
	index   uint64
	term    uint64
	command *Command
}

func (logentry *LogEntry) Error() string {
	panic("implement me")
}

func (entry *LogEntry) GetCommand() *Command {
	return entry.command
}

func (entry *LogEntry) GetTerm() uint64 {
	return entry.term
}

func (logentry *LogEntry) LogEntry(entryBuilder LongEntryBuilder) *LogEntry {
	return &LogEntry{index: entryBuilder.term, term: entryBuilder.term, command: &entryBuilder.command}
}

type LongEntryBuilder struct {
	index   uint64
	term    uint64
	command Command
}

func (l *LongEntryBuilder) Index(val uint64) *LongEntryBuilder {
	l.index = val
	return l
}

func (l *LongEntryBuilder) Term(val uint64) *LongEntryBuilder {
	l.term = val
	return l
}

func (l *LongEntryBuilder) Command(c Command) *LongEntryBuilder {
	l.command = c
	return l
}

func build(entryBuilder LongEntryBuilder) *LogEntry {
	l := new(LogEntry)
	l.LogEntry(entryBuilder)
	return l
}
