package impl

import (
	"context"
	"fmt"
	"htl/myraft.com/channs"
	command "htl/myraft.com/command"
	"htl/myraft.com/entry"
	vote "htl/myraft.com/entry/RvoteParam"
	"htl/myraft.com/inter"
	"htl/myraft.com/memchange"
	raft_client "htl/myraft.com/raft.client"
	"htl/myraft.com/rpc"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Instance *DefaultNode
)

type DefaultNode struct {
	//选举时间间隔
	electionSpanTime int64
	//上一次选举时间
	preElectionTime time.Time
	//上一次心跳时间
	preHeartBeatTime time.Time
	//心跳间隔基数
	heartBeatTick        int64
	status               int
	peerset              command.PeerSet
	clustememshipchanges inter.ClusterMembershipChanges
	//sync data chan
	channs.SyncData
	currentTerm int64
	//当前获得选票的候选人id
	voteFor string
	/** 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号 */
	LogModule inter.LogModule
	//已经被提交的日志最大索引值
	commitIndex int64
	//最后被应用到状态机的日志条目（初始化0 持续递增）
	lastApplied int64
	//对于每个服务器需要发送给它的下一条日志的索引号，初始化的时候领导人的值需要+1
	nextIndexes map[command.Peer]int64
	//对于follower服务器已经复制给他的日志最大索引值
	matchIndexes map[command.Peer]int64
	//node节点是否已经启动
	started bool
	//节点配置
	config command.NodeConfig
	//rpc server
	RpcServer rpc.RpcServer
	//rpc client
	RpcClient rpc.RpcClient
	//状态机
	stateMachine inter.StateMachine
	//一致性模块
	consensus inter.Consenus
	//private HeartBeatTask heartBeatTask = new HeartBeatTask();
	//private ElectionTask electionTask = new ElectionTask();
	//private ReplicationFailQueueConsumer replicationFailQueueConsumer = new ReplicationFailQueueConsumer();
	//
	//private LinkedBlockingQueue<ReplicationFailModel> replicationFailQueue = new LinkedBlockingQueue<>(2048);
	//
	//
	//inter.Node
	raftchannelpool  channs.RaftChannelPool
	SynchronizedLock sync.Mutex
}

func GetInstance() *DefaultNode {
	if Instance != nil {
		return Instance
	} else {
		return &DefaultNode{}
	}
}

func (node *DefaultNode) SetNodeConfig(config command.NodeConfig) {
	node.config = config
}

func (node *DefaultNode) GetPeerSet() *command.PeerSet {
	return node.GetPeerSet()
}

func (node *DefaultNode) GetNodeStatus() int {
	return node.status
}

func (node *DefaultNode) GetNextIndexes() map[command.Peer]int64 {
	return node.nextIndexes
}

func (node *DefaultNode) GetMatchedIndexes() map[command.Peer]int64 {
	return node.matchIndexes
}

func (node *DefaultNode) Init() {

	if node.started {
		return
	}
	node.SynchronizedLock.Lock()
	defer node.SynchronizedLock.Unlock()

	stopchan := make(chan struct{}, 0)
	node.RpcServer.Start(stopchan)
	node.consensus = NewDefaultConsenus()
	node.LogModule = NewDefaultLogModule()
	node.stateMachine = NewDefaultStateMachine()

	node.newElectionScheduler(6 * time.Second)
	node.newHeartBeatScheduler(500 * time.Millisecond)
	logentry := node.LogModule.GetLast()
	if logentry != nil {
		node.currentTerm = logentry.GetTerm()
	}

	node.started = true
	fmt.Printf("raft server start success selfid %s", node.peerset.GetSelf())
}

func (node *DefaultNode) HandlerAppendEntries(param *entry.AppendEntryParam) *entry.AppendEntryResult {
	if param != nil {
		fmt.Printf("node receive node %s append entry, entry content = %s", param.GetLeaderId(), param.GetEntries())
	}
	result := node.consensus.AppendEntries(param)
	return result
}

func (node *DefaultNode) AddPeers(newPeer command.Peer) *memchange.ClusterMemberChageResult {
	return node.clustememshipchanges.AddPeer(&newPeer)
}

func (node *DefaultNode) RemovePeers(oldPeer command.Peer) *memchange.ClusterMemberChageResult {
	return node.clustememshipchanges.RemovePeer(&oldPeer)
}

func (node *DefaultNode) HandlerRequestVote(param vote.RvoteParam) *entry.RvoteResult {
	fmt.Printf("handlerRequestVote will be invoke, param info : %s", param)
	return nil
}

func (node *DefaultNode) newHeartBeatScheduler(duration time.Duration) {
	tick := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-tick.C:
				node.HeartBeat()
			}
		}
	}()
}

func (node *DefaultNode) newElectionScheduler(duration time.Duration) {
	tick := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-tick.C:
				node.Election()
			}
		}
	}()
}

//raft节点的心跳实现
func (node *DefaultNode) HeartBeat() {
	if node.status != command.LEADER {
		return
	}

	currentTime := time.Now()

	if int64(currentTime.Sub(node.preHeartBeatTime)) < node.heartBeatTick {
		return
	}

	node.preHeartBeatTime = time.Now()
	fmt.Println("nextindex %s", node.nextIndexes[*node.peerset.GetSelf()])

	for x := node.peerset.GetPeersWithOutSelf().Front(); x.Value != nil; x = x.Next() {
		p := x.Value.(command.Peer)
		//设置发送给其他节点的心跳信息
		builder := entry.Newbuilder().Entries(nil).ServerId(p.GetAddr()).LeaderId(p.GetAddr())
		param := entry.NewAppendEntryParam(builder)
		//初始化请求
		request := new(rpc.Request)
		request.SetCmd(rpc.A_ENTRIS)
		request.SetObj(param)
		request.SetUrl(p.GetAddr())
		channs.NewRunnable(func() error {
			response := node.RpcClient.Send(*request)
			result := response.GetResult().(entry.AppendEntryResult)
			term := result.GetTerm()
			if term > node.currentTerm {
				fmt.Errorf("self will become follower, he's term : {}, my term : {}", term, node.currentTerm)
				node.currentTerm = term
				node.voteFor = ""
				node.status = command.FOLLOWER
			}
			return nil
		})
		node.raftchannelpool.TaskPool = append(node.raftchannelpool.TaskPool)
	}

}

//raft节点投票选举
func (node *DefaultNode) Election() {

	task := &channs.ChannelTask{}
	task.Runnable.Runnablefunc = func() error {
		//如果节点状态是leader则返回
		if node.status == command.LEADER {
			return nil
		}
		current := time.Now()
		node.electionSpanTime = node.electionSpanTime + rand.Int63n(50)
		if current.After(node.preHeartBeatTime.Add(time.Duration(int64(time.Millisecond) * node.electionSpanTime))) {
			return nil
		}
		node.status = command.CANDIDATE
		fmt.Println("warnning:node {} will become CANDIDATE and start election leader, current term : [{}], LastEntry : [{}] ")
		node.preElectionTime = node.preElectionTime.
			Add(time.Duration(int64(time.Millisecond) + 15*int64(time.Millisecond)))
		node.currentTerm = node.currentTerm + 1
		//先把term+1
		//推荐自己
		node.voteFor = node.peerset.GetSelf().GetAddr()
		peers := node.peerset.GetPeersWithOutSelf()
		fmt.Println("peers size:", peers.Len())
		futureArray := make([]*channs.ChannelTask, 0)
		for i := peers.Front(); i != nil; i = i.Next() {
			peer := i.Value.(command.Peer)
			t := &channs.ChannelTask{
				Callable: channs.NewCallable(func() (interface{}, error) {
					lastTerm := int64(0)
					lastLogEntry := node.LogModule.GetLast()
					//获取最后一个日志的任期
					if lastLogEntry != nil {
						lastTerm = int64(lastLogEntry.GetTerm())
					}
					//创建投票参数以及请求
					voteParam := vote.NewRvoteParam().SetTerm(node.currentTerm).SetCandidateId(node.peerset.GetSelf().GetAddr()).
						SetLastLogIndex(node.LogModule.GetLastIndex()).SetLastLogTerm(lastTerm)
					request := rpc.NewRequest().SetCmd(rpc.R_VOTE).SetObj(voteParam).SetUrl(peer.GetAddr())
					response := node.RpcClient.Send(*request)
					return response, nil
				}, context.Background()),
			}
			futureArray = append(futureArray, t)
			var succeed int64
			wg := &sync.WaitGroup{}
			wg.Add(len(futureArray))
			for _, futrue := range futureArray {
				//异步执行节点间投票
				go func() {
					response := node.raftchannelpool.Get(*futrue)
					voteResult := response.GetResult().(entry.RvoteResult)
					//如果集群的节点投票给我，则原子计数器+1
					if voteResult.VoteGranted() && response.Err == nil {
						atomic.AddInt64(&succeed, 1)
					} else {
						resTerm := response.GetResult().(*entry.RvoteResult).Term()
						if resTerm >= node.currentTerm {
							node.currentTerm = resTerm
						}
					}
					wg.Done()
				}()
			}
			wg.Wait()
			fmt.Printf("node %s maybe become leader , success count = {} , status : {}", node.peerset.GetSelf(), succeed, node.status)
			// 如果投票期间,有其他服务器发送 appendEntry , 就可能从CANDIDATE变成 FOLLOWER ,这时,应该停止.
			if node.status == command.FOLLOWER {
				return nil
			}

			if int(succeed) >= peers.Len()/2 {
				fmt.Printf("warning: node %s become leader", node.peerset.GetSelf())
				node.status = command.LEADER
				node.peerset.SetLeader(node.peerset.GetSelf())
				node.voteFor = ""
				node.becomeLeaderToDoThing()
			} else {
				node.voteFor = ""
			}

		}
		//把任务批量插入到channel pool中
		return nil
	}

}

func (node *DefaultNode) becomeLeaderToDoThing() {
	node.nextIndexes = make(map[command.Peer]int64)
	node.matchIndexes = make(map[command.Peer]int64)

	for peer := node.peerset.GetPeersWithOutSelf().Front(); peer != nil; peer = peer.Next() {
		node.nextIndexes[peer.Value.(command.Peer)] = (node.LogModule.GetLastIndex() + 1)
		node.matchIndexes[peer.Value.(command.Peer)] = 0
	}
}

func (node *DefaultNode) handlerRequestVote(vote *vote.RvoteParam) *entry.RvoteResult {
	fmt.Println("handlerRequestVote will be invoke, param info : {}", vote)
	response := node.consensus.RequestVote(vote)
	return response
}
func (node *DefaultNode) handlerAppendEntries(param *entry.AppendEntryParam) *entry.AppendEntryResult {
	fmt.Println("handlerAppendEntries will be invoke, param info : {}", param)
	response := node.consensus.AppendEntries(param)
	return response
}

func (node *DefaultNode) redirect(request raft_client.ClientKVReq) *raft_client.ClientKVAck {
	req := rpc.NewRequest().SetUrl(node.peerset.GetLeader().GetAddr()).SetObj(request).SetCmd(rpc.CLIENT_REQ)
	response := node.RpcClient.Send(*req)
	return response.GetResult().(*raft_client.ClientKVAck)
}

/**
 * 客户端的每一个请求都包含一条被复制状态机执行的指令。
 * 领导人把这条指令作为一条新的日志条目附加到日志中去，然后并行的发起附加条目 RPCs 给其他的服务器，让他们复制这条日志条目。
 * 当这条日志条目被安全的复制（下面会介绍），领导人会应用这条日志条目到它的状态机中然后把执行的结果返回给客户端。
 * 如果跟随者崩溃或者运行缓慢，再或者网络丢包，
 *  领导人会不断的重复尝试附加日志条目 RPCs （尽管已经回复了客户端）直到所有的跟随者都最终存储了所有的日志条目。
 */
func (node *DefaultNode) HandlerClientRequest(request raft_client.ClientKVReq) *raft_client.ClientKVAck {
	fmt.Println("handlerClientRequest handler {} operation,  and key : [{}], value : [{}]", request)

	//如果当前节点不是leader则把请求转发到leader处理
	if node.status != command.LEADER {
		fmt.Printf(" warning I not am leader , only invoke redirect method, leader addr : %s, my addr : %s",
			node.peerset.GetLeader().GetAddr(), node.peerset.GetSelf().GetAddr())
		return node.redirect(request)
	}

	//如果用户的请求是获取raft 日志则直接从状态机获取返回
	if request.GetReqType() == raft_client.GET {
		logEntry := node.stateMachine.Get(request.GetKey())
		if logEntry != nil {
			return raft_client.NewClientKVAck().Ojbect(logEntry.GetCommand())
		}
		return raft_client.NewClientKVAck().Ojbect(nil)
	}

	//创建日志
	logentryBuilder := &entry.LogEntryBuilder{}
	logentry := logentryBuilder.Term(node.currentTerm).
		Command(entry.NewCommand().SetKey(request.GetKey()).SetValue(request.GetValue())).LogEntryBuild()
	// 预提交到本地日志,
	node.LogModule.Write(logentry)
	fmt.Printf("write logModule success, logEntry info : %s, log index : %v", logentry, logentry.GetIndex())
	writeTasksArrys := make([]*channs.ChannelTask, 0)
	var count int = 0
	for peer := node.peerset.GetPeersWithOutSelf().Front(); peer != nil; peer = peer.Next() {
		writeTasksArrys = append(writeTasksArrys, node.Replication(peer.Value.(command.Peer), logentry))
		count++
	}

	node.raftchannelpool.TaskPool = append(node.raftchannelpool.TaskPool, writeTasksArrys...)
	var success int64 = 0
	//异步返回结果统计执行成功的次数
	for _, task := range writeTasksArrys {
		result := node.raftchannelpool.Get(*task)
		if result.GetResult().(bool) {
			atomic.AddInt64(&success, 1)
		}
	}

	// 如果存在一个满足N > commitIndex的 N，并且大多数的matchIndex[i] ≥ N成立，
	// 并且log[N].term == currentTerm成立，那么令 commitIndex 等于这个 N （5.3 和 5.4 节）

	matchIndexesList := make([]int64, 0)
	for _, value := range node.matchIndexes {
		matchIndexesList = append(matchIndexesList, value)
	}

	middle := 0
	if len(matchIndexesList) >= 2 {
		middle = len(matchIndexesList) / 2
	}

	num := matchIndexesList[middle]

	if num > node.commitIndex {
		log := node.LogModule.Read(num)
		if log != nil && log.GetTerm() == node.currentTerm {
			node.commitIndex = num
		}
	}

	if int(success) >= count/2 {
		node.commitIndex = logentry.GetIndex()
		node.stateMachine.Apply(logentry)
		node.lastApplied = node.commitIndex
		fmt.Printf("success apply local state machine,  logEntry info : %s", logentry)
		return raft_client.NewClientKVAck().Ok()
	} else {
		node.LogModule.RemoveOnstartIndex(logentry.GetIndex())
		fmt.Printf("fail apply local state  machine,  logEntry info : %s", logentry)
		return raft_client.NewClientKVAck().Fail()
	}
}

//leader 把日志复制到其他node节点
func (node *DefaultNode) Replication(peer command.Peer, logentry *entry.LogEntry) *channs.ChannelTask {
	callable := channs.NewCallable(
		func() (interface{}, error) {
			var start, end time.Time = time.Now(), time.Now()
			for !end.After(start.Add(time.Duration(time.Duration(20 * time.Millisecond)))) {

				param := entry.Newbuilder().Term(node.currentTerm).
					ServerId(peer.GetAddr()).
					LeaderId(node.peerset.GetSelf().GetAddr()).Leadercommit(node.commitIndex)

				// 以我这边为准, 这个行为通常是成为 leader 后,首次进行 RPC 才有意义.
				nextindex := node.nextIndexes[peer]
				listLogEntry := make([]entry.LogEntry, 0)
				if logentry.GetIndex() >= nextindex {
					for i := nextindex; i <= logentry.GetIndex(); i++ {
						l := node.LogModule.Read(i)
						if l != nil {
							listLogEntry = append(listLogEntry, *l)
						}
					}
				} else {
					listLogEntry = append(listLogEntry, *logentry)
				}

				preLog := listLogEntry[0]
				param.PreLogTerm(preLog.GetTerm())
				param.PreLogIndex(preLog.GetIndex())
				param.Entries(listLogEntry)

				AppendEntryParam := entry.NewAppendEntryParam(param)

				request := rpc.NewRequest().SetCmd(rpc.A_ENTRIS).SetObj(AppendEntryParam).SetUrl(peer.GetAddr())
				response := node.RpcClient.Send(*request)
				if response == nil {
					return false, nil
				}

				result := response.GetResult().(*entry.AppendEntryResult)
				if result != nil && result.GetStatus() {
					fmt.Printf("append follower entry success , follower=%s, entry=%s", peer, AppendEntryParam.GetEntries())
					node.nextIndexes[peer] = logentry.GetIndex() + 1
					node.matchIndexes[peer] = logentry.GetIndex()
					return true, nil
				} else if result != nil {
					if result.GetTerm() > node.currentTerm {

						fmt.Printf("follower %s term %d than more self, and my term = [{}], so, I will become follower",
							peer, result.GetTerm(), node.currentTerm)
						node.currentTerm = result.GetTerm()
						// 认怂, 变成跟随者
						node.status = command.FOLLOWER
						return false, nil
					} else {
						node.nextIndexes[peer] = nextindex - 1
					}
				}
				end = time.Now()
			}
			return false, nil
		},
		context.Background(),
	)

	c := &channs.ChannelTask{Runnable: nil, Callable: callable}
	return c
}
