package lightswitch

import (
	"sync"
)

var mutex = &sync.Mutex{}
var counter = 0

func Lock(m sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	
	counter++
	if counter == 1 {
		m.Lock()
	}
}

func Unlock(m sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	
	counter--
	if counter == 0 {
		m.Unlock()
	}
}
