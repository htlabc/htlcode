package inter

type LifeCycle interface {
	init() error
	destory() error
}
