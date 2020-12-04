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
	//对于每个服务器需要发送给它的下一条日志的索引号，初始化领导人需要+1
	nextIndexes map[command.Peer]uint64
	//对于follower服务器已经复制给他的日志最大索引值
	matchIndexes map[command.Peer]uint64
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
	raftchannelqueue channs.RaftChannelPool
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

func (h *HeartBeatTask) HeartBeatTask() {
	if h.node.status != command.LEADER {
		return
	}

	currentTime := time.Now()

	if int64(currentTime.Sub(h.node.preHeartBeatTime)) < h.node.heartBeatTick {
		return
	}

	h.node.preHeartBeatTime = time.Now()
	fmt.Println("nextindex")

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
		channs.RaftChannelPools.AddTasks(func() error {
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

	}

}

func (h *HeartBeatTask) ElectionTask() {

	task := &channs.ChannelTask{}
	task.Runnable.Runnablefunc = func() error {
		//如果节点状态是leader则返回
		if h.node.status == command.LEADER {
			return nil
		}
		current := time.Now()
		h.node.electionSpanTime = h.node.electionSpanTime + rand.Int63n(50)
		if current.After(h.node.preHeartBeatTime.Add(time.Duration(int64(time.Millisecond) * h.node.electionSpanTime))) {
			return nil
		}
		h.node.status = command.CANDIDATE
		fmt.Println("warnning:node {} will become CANDIDATE and start election leader, current term : [{}], LastEntry : [{}] ")
		h.node.preElectionTime = h.node.preElectionTime.
			Add(time.Duration(int64(time.Millisecond) + 15*int64(time.Millisecond)))
		h.node.currentTerm = h.node.currentTerm + 1
		//先把term+1
		//推荐自己
		h.node.voteFor = h.node.peerset.getAddr()
		peers := h.node.peerset.GetPeersWithOutSelf()
		fmt.Println("peers size:", peers.Len())
		futureArray := make([]*channs.ChannelTask, 0)
		for i := peers.Front(); i != nil; i = i.Next() {
			peer := i.Value.(command.Peer)
			t := &channs.ChannelTask{
				Callable: channs.NewCallable(func() (interface{}, error) {

					lastTerm := int64(0)
					lastLogEntry := h.node.logModule.getLast()
					//获取最后一个日志的任期
					if lastLogEntry != nil {
						lastTerm = lastLogEntry.getTerm()
					}

					voteParam := vote.NewRvoteParam().SetTerm(h.node.currentTerm).SetCandidateId(h.node.peerset.getSelf().getAddr()).
						SetLastLogIndex(h.node.logModule.getLastIndex()).SetLastLogTerm(lastTerm)

					return "", nil
				}, context.Background()),
			}
			futureArray = append(futureArray, t)
		}
		//把任务批量插入到channel pool中
		h.node.raftchannelqueue.TaskPool = append(h.node.raftchannelqueue.TaskPool, futureArray...)

		return nil
	}
	h.node.raftchannelqueue.TaskPool = append(h.node.raftchannelqueue.TaskPool, task)

}
