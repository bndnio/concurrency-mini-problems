package main

import (
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
	
	for i:=0; i<100; i++ {
		wg.Add(1)
		go prod(i)
	}
	wg.Add(1)
	go cons()
	wg.Wait()

	if *memprofile != "" {
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
