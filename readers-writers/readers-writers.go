package main

import (
	// "fmt"
	"sync"
	"errors"
	"math/rand"
)

var roomEmpty = &sync.Mutex{}
var mutex = &sync.Mutex{}
var readerCount = 0
var wg sync.WaitGroup
var ds = make([]int, 0)

func reader() (int, error) {
	defer wg.Done()
	
	mutex.Lock()
	if readerCount == 0 {
		roomEmpty.Lock()
	}
	readerCount++
	mutex.Unlock()

	// fmt.Println("read: ", out)
	if len(ds) == 0 {return -1, errors.New("empty queue") }
	out := ds[rand.Intn(len(ds))]

	mutex.Lock()
	readerCount--
	if readerCount == 0 {
		roomEmpty.Unlock()
	}
	mutex.Unlock()
	return out, nil
}

func writer(num int) {
	roomEmpty.Lock()
	
	// *** START CRITICAL SECTION ***
	ds = append(ds, num)
	// fmt.Println("wrote: ", num)
	// *** END CRITICAL SECTION ***

	wg.Done()
	roomEmpty.Unlock()
}

func main() {
	for i:=0; i<110; i++ {
		wg.Add(1)
		if i % 10 == 0 {
			go writer(i)
		} else {
			go reader()
		}
	}
	wg.Wait()
}