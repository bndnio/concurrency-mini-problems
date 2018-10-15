package main

import (
	// "fmt"
	"sync"
	"container/list"
	"lightswitch_channel"
)

var l *list.List = nil
var wg sync.WaitGroup

var insertMutex = make(chan bool, 1)
var noSearcher = make(chan bool, 1)
var noInserter = make(chan bool, 1)
var searchSwitch = lightswitch_channel.New()
var insertSwitch = lightswitch_channel.New()

func insert(value int) {
	defer wg.Done()
	searchSwitch.Lock(noSearcher)
	l.PushBack(value)
	searchSwitch.Unlock(noSearcher)
}

func find(value int) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return e
		}
	}
	return nil
}

func search(value int) *list.Element {
	defer wg.Done()
	insertSwitch.Lock(noInserter)
	<- insertMutex
	insertMutex <- true
	insertSwitch.Unlock(noInserter)

	return find(value)
}

func delete(value int) {
	defer wg.Done()
	<- noSearcher
	<- noInserter
	l.Remove(find(value))
	noInserter <- true
	noSearcher <- true
}

func print() {
	for e := l.Front(); e != nil; e = e.Next() {
		// fmt.Print(e.Value, " ")
	}
	// fmt.Print("\n")
}

func main() {
	insertMutex <- true
	noSearcher <- true
	noInserter <- true
	l = list.New()
	for i:=0; i<100; i++ {
		if i%2 == 0 {
			wg.Add(1)
			insert(i)
		} else if i%3 == 0 {
			wg.Add(10)
			for j:=0; j<10; j++ {
				search(i-1)
			}
		} else {
			wg.Add(1)
			delete(i-1)
		}
		print()
	}

	wg.Wait()
}
