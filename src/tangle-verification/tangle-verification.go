package main

import (
	"fmt"
	"sync"
	"tangle_node"
)

var wg sync.WaitGroup

var queue = make([]*tangle_node.Node, 0)
var queueWriteMutex = make(chan bool, 1)

func addWork() {
	defer wg.Done()
	newNode := tangle_node.New()

	var isComplete int
	for verified, i := 0, 0; verified < 2; i++ {
		var comm = make(chan bool, 1)
		var cb = make(chan bool)
		fmt.Println(len(queue))
		go queue[i%len(queue)].Verify(newNode, comm, cb)
		// ** node evaluation work would go here **
		comm <- true
		didVerify := <- cb
		if didVerify { verified++ }
	}
	<- queueWriteMutex
	queue = append(queue, newNode)
	if isComplete {
		queue = append(queue[:i%len(queue)], queue[((i%len(queue))+1):]...)
	}
	queueWriteMutex <- true
}

// Execution code
func main() {
	queue = append(queue, tangle_node.New())
	// head := queue[0]
	queue = append(queue, tangle_node.New())
	queueWriteMutex <- true

	for i:=0; i<80; i++ {
		wg.Add(1)
		go addWork()
	}

	wg.Wait()
}
