package tangle_node

import "fmt"
import "sync"

type Node struct {
	verified bool
	approvalCount int
	approvalMutex chan bool
	workingMutex chan bool
	outbound []*Node
}

var nd *Node

func New() *Node {
	nd = &Node{
		verified: false, 
		approvalCount: 0, 
		approvalMutex: make(chan bool, 1), 
		workingMutex: make(chan bool, 2),
		outbound: make([]*Node, 0),
	}
	nd.approvalMutex <- true
	return nd
}

func inSlice(ele *Node, slc []*Node) bool {
	for i := range slc {
		if slc[i] == ele {
			return true
		}
	}
	return false
}

func (nd *Node) Print(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(*nd)
	for i := range nd.outbound {
		wg.Add(1)
		go nd.outbound[i].Print(wg)
	}
}

func (nd *Node) Verify(addr *Node, approved <-chan bool, cb chan<- bool) {
	// check if Node already approved
	<- nd.approvalMutex
	if nd.approvalCount >= 2 || inSlice(addr, nd.outbound) == true {
		cb <- false
		nd.approvalMutex <- true
		return
	}
	nd.approvalMutex <- true

	// allow requester to do work
	nd.workingMutex <- true
	var verified = <- approved
	if verified == true {
		// check for approval account again
		<- nd.approvalMutex
		// if full return false
		// otherwise increase approval count
		if nd.approvalCount >= 2 || inSlice(addr, nd.outbound) == true  {
			fmt.Println("mid fail")
			cb <- false
			nd.approvalMutex <- true 
			<- nd.workingMutex
			return
		} else { 
			nd.approvalCount++
		}
		// and mark verified if possible
		if nd.approvalCount == 2 {
			nd.verified = true
		}
		// and set outbound link to Node
		nd.outbound = append(nd.outbound, addr)
		nd.approvalMutex <- true 
	}
	<- nd.workingMutex
	cb <- true
}