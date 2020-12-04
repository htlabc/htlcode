package entry

//添加新日志
type AppendEntryParam struct {
	//领导人id
	leaderId string
	//新日志的任期号
	term uint64
	//服务器id
	serverId string
	//新日志的上一条日志
	preLogIndex uint64
	//新日志的上一条日志的任期号
	preLogTerm uint64
	//raft 日志集合
	entries []LogEntry
	//领导人已经提交的日志索引值
	leaderCommit uint64
}

func NewAppendEntryParam(b *builder) *AppendEntryParam {
	p := &AppendEntryParam{}
	p.term = b.term
	p.leaderId = b.leaderid
	p.leaderCommit = b.leaderCommit
	p.preLogIndex = b.preLogIndex
	p.preLogTerm = b.preLogTerm
	p.entries = b.entries
	p.serverId = b.serverId
	return p
}
func Newbuilder() *builder {
	b := new(builder)
	return b
}

type builder struct {
	//private long term;
	//private String serverId;
	//private String leaderId;
	//private long prevLogIndex;
	//private long preLogTerm;
	//private LogEntry[] entries;
	//private long leaderCommit;

	term         uint64
	serverId     string
	leaderid     string
	preLogIndex  uint64
	preLogTerm   uint64
	leaderCommit uint64
	entries      []LogEntry
}

func (b *builder) Term(val uint64) *builder {
	b.term = val
	return b
}

func (b *builder) LeaderId(val string) *builder {
	b.leaderid = val
	return b
}

func (b *builder) Leadercommit(val uint64) *builder {
	b.leaderCommit = val
	return b
}

func (b *builder) ServerId(val string) *builder {
	b.serverId = val
	return b
}

func (b *builder) PreLogIndex(val uint64) *builder {
	b.preLogIndex = val
	return b
}

func (b *builder) LogTerm(val uint64) *builder {
	b.preLogTerm = val
	return b
}

func (b *builder) Entries(val []LogEntry) *builder {
	b.entries = val
	return b
}
