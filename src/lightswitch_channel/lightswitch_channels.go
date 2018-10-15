package lightswitch_channel

// import "fmt"

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
	<- (ls.mutex)
	
	ls.counter -= 1
	if ls.counter == 0 {
		c <- true
	}

	ls.mutex <- true
}
