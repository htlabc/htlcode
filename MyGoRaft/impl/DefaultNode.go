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
	logModule inter.LogModule
	//已经被提交的日志最大索引值
	commitIndex uint64
	//最后被应用到状态机的日志条目（初始化0 持续递增）
	lastApplied int
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
	raftchannelpool channs.RaftChannelPool
}

func (node *DefaultNode) HandlerAppendEntries(param entry.AppendEntryParam) *entry.AppendEntryResult {
	return nil
}

func (node *DefaultNode) HandlerClientRequest(param raft_client.ClientKVReq) *raft_client.ClientKVAck {
	return nil
}

func (node *DefaultNode) AddPeers(newPeer command.Peer) memchange.ClusterMemberChageResult {
	return node.clustememshipchanges.AddPeer(newPeer)
}

func (node *DefaultNode) RemovePeers(oldPeer command.Peer) memchange.ClusterMemberChageResult {
	return node.clustememshipchanges.RemovePeer(oldPeer)
}

func (node *DefaultNode) HandlerRequestVote(param vote.RvoteParam) *entry.RvoteResult {
	fmt.Printf("handlerRequestVote will be invoke, param info : %s", param)
	return nil
}

type HeartBeatTask struct {
	tick time.Ticker
	node DefaultNode
}

//raft节点的心跳实现
func (h *HeartBeatTask) HeartBeatTask() {
	if h.node.status != command.LEADER {
		return
	}

	currentTime := time.Now()

	if int64(currentTime.Sub(h.node.preHeartBeatTime)) < h.node.heartBeatTick {
		return
	}

	h.node.preHeartBeatTime = time.Now()
	fmt.Println("nextindex %s", h.node.nextIndexes[*h.node.peerset.GetSelf()])

	for x := h.node.peerset.GetPeersWithOutSelf().Front(); x.Value != nil; x.Next() {
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
			response := h.node.RpcClient.Send(*request)
			result := response.GetResult().(entry.AppendEntryResult)
			term := result.GetTerm()
			if term > h.node.currentTerm {
				fmt.Errorf("self will become follower, he's term : {}, my term : {}", term, h.node.currentTerm)
				h.node.currentTerm = term
				h.node.voteFor = ""
				h.node.status = command.FOLLOWER
			}
			return nil
		})
		h.node.raftchannelpool.TaskPool = append(h.node.raftchannelpool.TaskPool)
	}

}

//raft节点投票选举
func (node *DefaultNode) ElectionTask() {

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
					lastLogEntry := node.logModule.GetLast()
					//获取最后一个日志的任期
					if lastLogEntry != nil {
						lastTerm = int64(lastLogEntry.GetTerm())
					}
					//创建投票参数以及请求
					voteParam := vote.NewRvoteParam().SetTerm(node.currentTerm).SetCandidateId(node.peerset.GetSelf().GetAddr()).
						SetLastLogIndex(node.logModule.GetLastIndex()).SetLastLogTerm(lastTerm)
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

	for peer := node.peerset.GetPeersWithOutSelf().Front(); peer != nil; peer.Next() {
		node.nextIndexes[peer.Value.(command.Peer)] = (node.logModule.GetLastIndex() + 1)
		node.matchIndexes[peer.Value.(command.Peer)] = 0
	}
}

func (node *DefaultNode) handlerRequestVote(vote vote.RvoteParam) *entry.RvoteResult {
	fmt.Println("handlerRequestVote will be invoke, param info : {}", vote)
	response := node.consensus.RequestVote(vote)
	return &response
}
func (node *DefaultNode) handlerAppendEntries(param entry.AppendEntryParam) *entry.AppendEntryResult {
	fmt.Println("handlerAppendEntries will be invoke, param info : {}", param)
	response := node.consensus.AppendEntries(param)
	return &response
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
func (node *DefaultNode) handlerClientRequest(request raft_client.ClientKVReq) *raft_client.ClientKVAck {
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

	node.logModule.Write(logentry)

	return nil
}
