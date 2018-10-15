/*
 * Note:
 * This code is interpretted, modified, and applied from the
 * pseudo code provided in the "Little Book of Semaphores"
 * written by Allen B. Downey (version 2.2.1)
 */

package lightswitch_mutex

import "sync"

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
	ls.counter += 1
	if ls.counter == 1 {
		m.Lock()
	}
	ls.mutex.Unlock()
}

func (ls *lightswitch) Unlock(m *sync.Mutex) {
	ls.mutex.Lock()
	ls.counter -= 1
	if ls.counter == 0 {
		m.Unlock()
	}
	ls.mutex.Unlock()
}
