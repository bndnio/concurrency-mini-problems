package main

import (
	"fmt"
	"sync"
	"math/rand"
)

var roomEmpty = &sync.Mutex{}
var mutex = &sync.Mutex{}
var readerCount = 0
var wg sync.WaitGroup
var ds = make([]int, 0)

func reader() int {
	mutex.Lock()
	if readerCount == 0 {
		roomEmpty.Lock()
	}
	readerCount++
	mutex.Unlock()

	out := ds[rand.Intn(len(ds))]
	fmt.Println("read: ", out)

	mutex.Lock()
	readerCount--
	if readerCount == 0 {
		roomEmpty.Unlock()
	}
	mutex.Unlock()
	defer wg.Done()
	return out
}

func writer(num int) {
	roomEmpty.Lock()
	
	// *** START CRITICAL SECTION ***
	ds = append(ds, num)
	fmt.Println("wrote: ", num)
	// *** END CRITICAL SECTION ***

	wg.Done()
	roomEmpty.Unlock()
}

func main() {
	for i:=0; i<10; i++ {
		wg.Add(1)
		go writer(i)
	}
	for i:=0; i<100; i++ {
		wg.Add(1)
		go reader()
	}
	wg.Wait()
}