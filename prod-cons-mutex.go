package main

import (
	"fmt"
	"sync"
	"errors"
	"time"
)

var mutex = &sync.Mutex{}
var wg sync.WaitGroup
var queue = make([]int, 0)

func enqueue(num int) {
	mutex.Lock()
	defer wg.Done()
	defer mutex.Unlock()

	// *** START CRITICAL SECTION ***
	queue = append(queue, num)
	// *** END CRITICAL SECTION ***
}

func dequeue() (int, error) {
	mutex.Lock()
	defer wg.Done()
	defer mutex.Unlock()

	// *** START CRITICAL SECTION ***
	if len(queue) == 0 { return -1, errors.New("empty queue") }
	out := queue[0]
	queue = queue[1:]
	// *** END CRITICAL SECTION ***

	return out, nil
}

func prod(i int) {
	for j:=0; j<10; j++ {
		wg.Add(1)
		go enqueue(i*1000+j)
	}
	wg.Done()
}

func cons() {
	failedAttempts := 0
	for {
		wg.Add(1)

		out, err := dequeue()
		if (err == nil) {
			failedAttempts = 0
			fmt.Println("consumed: ", out)
		} else {
			if (failedAttempts == 10) {
				fmt.Println("goodbye")
				break
			} else {
				failedAttempts++
				time.Sleep(100*time.Millisecond)
			}
		}
	}
	wg.Done()
}

func main() {
	for i:=0; i<10; i++ {
		wg.Add(1)
		go prod(i)
	}
	for i:=0; i<10; i++ {
		wg.Add(1)
		go cons()
	}
	wg.Wait()
}
