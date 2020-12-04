package channs

type SyncData struct {
	readchan  chan interface{}
	writechan chan interface{}
}

func (s *SyncData) syncGetData() {

	var value interface{}
	go func() {
		select {
		case value = <-s.writechan:
		case s.readchan <- value:
		}
	}()

}

func NewSyncData() *SyncData {
	s := new(SyncData)
	s.writechan, s.readchan = make(chan interface{}), make(chan interface{})
	return s
}
