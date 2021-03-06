package main

import (
	// "fmt"
	"sync"
	"errors"
	"time"
	"flag"
	"os"
	"log"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
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
	for j:=0; j<1000; j++ {
		wg.Add(1)
		go enqueue(i*1000+j)
	}
	wg.Done()
}

func cons() {
	failedAttempts := 0
	for {
		wg.Add(1)

		_, err := dequeue()
		if (err == nil) {
			failedAttempts = 0
			// fmt.Println("consumed: ", out)
			wg.Add(1)
			go cons()
		} else {
			if (failedAttempts == 2) {
				// fmt.Println("goodbye")
				break
			} else {
				failedAttempts++
				time.Sleep(time.Duration(100*failedAttempts)*time.Millisecond)
			}
		}
	}
	wg.Done()
}

func main() {
	flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
	}

	for i:=0; i<1000; i++ {
		wg.Add(1)
		go prod(i)
	}
	wg.Add(1)
	go cons()
	wg.Wait()

	if *memprofile != "" {
		runtime.MemProfileRate = 1
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
    }
}
