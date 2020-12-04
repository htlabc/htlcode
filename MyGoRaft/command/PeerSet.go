package command

import (
	"container/list"
	"sync"
)

type PeerSet struct {
	list *list.List
	//volatile peer
	leader Peer
	//volatile
	self *Peer
}

var p *PeerSet
var once sync.Once

func GetInstance() *PeerSet {
	once.Do(
		func() {
			p = &PeerSet{}
			p.list = list.New()
		})
	return p
}

func (p *PeerSet) AddPeer(peer *Peer) {
	p.list.PushBack(peer)
}

func (p *PeerSet) RemovePeer(peer *Peer) {
	for e := p.list.Front(); e != nil; e = e.Next() {
		if e.Value == peer {
			p.list.Remove(e)
		}
	}
}

func (p *PeerSet) GetSelf() *Peer {
	return p.self
}

func (p *PeerSet) GetPeers() *list.List {
	return p.list
}

func (p *PeerSet) GetPeersWithOutSelf() *list.List {
	l := list.New()
	l = p.list
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == p.self {
			p.list.Remove(e)
		}
	}
	return l
}

func (p *PeerSet) GetLeader() Peer {
	return p.leader
}

func (p *PeerSet) SetLeader(value Peer) {
	p.leader = value
}
