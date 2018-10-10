package main

import (
	"fmt"
	"sync"
	"errors"
	"time"
)

var wg sync.WaitGroup
var queue = make(chan int, 100)

func enqueue(num int) {
	defer wg.Done()

	// *** START CRITICAL SECTION ***
	queue <- num
	// *** END CRITICAL SECTION ***
}

func dequeue() (int, error) {
	defer wg.Done()

	// *** START CRITICAL SECTION ***
	select {
	case out := <- queue: return out, nil
	default: return -1, errors.New("empty queue") 
	}
	// *** END CRITICAL SECTION ***
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
