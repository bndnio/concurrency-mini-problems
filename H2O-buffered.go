package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var hReady = make(chan bool, 2)
var oReady = make(chan bool, 2)

func H() {
	hReady <- true
	<- oReady
	fmt.Println("H through")
	wg.Done()
}

func O() {
	<- hReady
	<- hReady
	oReady <- true
	oReady <- true
	fmt.Println("O through")
	wg.Done()
}

func main() {
	for i:=0; i<20; i++ {
		wg.Add(1)
		go H()
	}
	for i:=0; i<10; i++ {
		wg.Add(1)
		go O()
	}
	wg.Wait()
}
