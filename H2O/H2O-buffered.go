package main

import (
	// "fmt"
	"sync"
	"flag"
	"os"
	"log"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var wg sync.WaitGroup
var hReady = make(chan bool, 2)
var oReady = make(chan bool, 2)

func H() {
	hReady <- true
	<- oReady
	// fmt.Println("H through")
	wg.Done()
}

func O() {
	<- hReady
	<- hReady
	oReady <- true
	oReady <- true
	// fmt.Println("O through")
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

	for i:=0; i<1000000; i++ {
		wg.Add(3)
		go H()
		go O()
		go H()
	}
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
