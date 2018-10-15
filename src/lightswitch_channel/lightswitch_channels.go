/*
 * Note:
 * This code is interpretted, modified, and applied from the
 * pseudo code provided in the "Little Book of Semaphores"
 * written by Allen B. Downey (version 2.2.1)
 */

package lightswitch_channel

type lightswitch struct {
	mutex chan bool
	counter int
}

var ls *lightswitch

func New() *lightswitch {
	ls = &lightswitch{make(chan bool, 1), 0}
	ls.mutex <- true
	return ls
}

func (ls *lightswitch) Lock(c chan bool) {
	<- ls.mutex
	ls.counter += 1
	if ls.counter == 1 {
		<- c
	}
	ls.mutex <- true
}

func (ls *lightswitch) Unlock(c chan bool) {
	<- ls.mutex
	ls.counter -= 1
	if ls.counter == 0 {
		c <- true
	}
	ls.mutex <- true
}
