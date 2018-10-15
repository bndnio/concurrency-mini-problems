package main

import (
	// "fmt"
	"sync"
	"errors"
	"math/rand"
	"flag"
	"os"
	"log"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
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
	if len(ds) == 0 {
		mutex.Lock()
		readerCount--
		if readerCount == 0 {
			roomEmpty.Unlock()
		}
		mutex.Unlock()
		return -1, errors.New("empty queue")
	}
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
	defer wg.Done()
	roomEmpty.Lock()
	
	// *** START CRITICAL SECTION ***
	ds = append(ds, num)
	// fmt.Println("wrote: ", num)
	// *** END CRITICAL SECTION ***

	roomEmpty.Unlock()
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

	for i:=0; i<1000000; i++ {
		wg.Add(1)
		if i % 10 == 0 {
			go writer(i)
		} else {
			go reader()
		}
	}
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