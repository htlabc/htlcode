package command

type Peer struct {
	addr string
}

func (p *Peer) GetAddr() string {
	return p.addr
}

func (p *Peer) toString() string {
	return "Peer { " + p.addr + "}"
}
