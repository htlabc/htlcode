package channs

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ChannelTask struct {
	taskid   string
	Runnable *Runnable
	Callable *Callable
}

type Runnable struct {
	Runnablefunc func() error
}

func NewCallable(f func() (interface{}, error), ctx context.Context) *Callable {
	ctx, cancel := context.WithCancel(ctx)
	return &Callable{
		callableFunc:   f,
		callbaleResult: &CallableResult{},
		cancel:         cancel,
		ctx:            ctx,
	}
}

func NewRunnable(f func() error) *Runnable {
	return &Runnable{
		Runnablefunc: f,
	}
}

type Callable struct {
	callableFunc   func() (interface{}, error)
	callbaleResult *CallableResult
	cancel         context.CancelFunc
	ctx            context.Context
}

type CallableResult struct {
	result interface{}
	err    error
}

func (c *ChannelTask) execute() error {
	if c.Runnable != nil {
		return c.Runnable.Runnablefunc()
	}

	if c.Callable != nil {
		result, err := c.Callable.callableFunc()
		c.Callable.callbaleResult.result = result
		c.Callable.callbaleResult.err = err
		c.Callable.ctx.Done()
		if c.Callable.callbaleResult.err != nil {
			return c.Callable.callbaleResult.err
		}
	}
	return nil

}

type RaftChannelPool struct {
	//协程池任务入口
	TaskPool             []*ChannelTask
	taskidMapChannelTask map[string]*ChannelTask
	channelTask          chan *ChannelTask
	//协程池任务数
	workerNum int
	//内部任务就绪队列
	jobsChannel chan *ChannelTask
	//
	lock sync.Mutex
}

func (p *RaftChannelPool) addTasksChannel() error {
	defer close(p.channelTask)
	for {
		for _, tsk := range p.TaskPool {
			if tsk.taskid == "" {
				tsk.taskid = GetRandomString(8)
			}
			p.taskidMapChannelTask[tsk.taskid] = tsk
			p.channelTask <- tsk
		}
	}
	return nil
}

func NewChannelPool(cap int) *RaftChannelPool {
	p := &RaftChannelPool{
		channelTask: make(chan *ChannelTask, cap),
		workerNum:   cap,
		jobsChannel: make(chan *ChannelTask, cap),
	}
	go p.addTasksChannel()
	return p
}

func (p *RaftChannelPool) worker(worker_id int) {
	for task := range p.jobsChannel {
		task.execute()
		fmt.Printf("workder id %d", worker_id)
	}
}

func (p *RaftChannelPool) Run() {
	//1,首先根据协程池的worker数量限定,开启固定数量的Worker,
	//  每一个Worker用一个Goroutine承载

	//3, 执行完毕需要关闭JobsChannel
	defer close(p.jobsChannel)
	//4, 执行完毕需要关闭EntryChannel
	defer close(p.channelTask)

	var num int
	if len(p.TaskPool) <= cap(p.channelTask) {
		num = len(p.TaskPool)
	} else {
		num = cap(p.channelTask)
	}

	for i := 0; i < num; i++ {
		go p.worker(i)
	}

	//2, 从EntryChannel协程池入口取外界传递过来的任务
	//   并且将任务送进JobsChannel中
	for task := range p.channelTask {
		p.jobsChannel <- task
	}

}

func (p *RaftChannelPool) Get(task ChannelTask) *CallableResult {
	if exist := p.taskidMapChannelTask[task.taskid]; exist != nil {
		exist.Callable.cancel()
		return exist.Callable.callbaleResult
	}

	return nil
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
