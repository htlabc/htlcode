package inter

type LifeCycle interface {
	Init() error
	Destory() error
}
