package lightswitch_mutex

import (
	"sync"
)

type lightswitch struct {
	mutex *sync.Mutex
	counter int
}

var ls *lightswitch

func New() *lightswitch {
	ls = &lightswitch{&sync.Mutex{}, 0}
	return ls
}

func (ls *lightswitch) Lock(m *sync.Mutex) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	
	ls.counter += 1
	if ls.counter == 1 {
		m.Lock()
	}
}

func (ls *lightswitch) Unlock(m *sync.Mutex) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	
	ls.counter -= 1
	if ls.counter == 0 {
		m.Unlock()
	}
}
