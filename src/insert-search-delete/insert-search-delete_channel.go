/*
 * Note:
 * This code is interpretted, modified, and applied from the
 * pseudo code provided in the "Little Book of Semaphores"
 * written by Allen B. Downey (version 2.2.1)
 */

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
	found := find(value)
	insertMutex <- true
	insertSwitch.Unlock(noInserter)

	return found
}

func delete(value int) {
	defer wg.Done()
	<- noSearcher
	<- noInserter
	toRemove := find(value)
	if toRemove != nil { l.Remove(toRemove) }
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
	for i:=0; i<1000000; i++ {
		if i%2 == 0 {
			wg.Add(1)
			go insert(i)
		} else if i%3 == 0 {
			wg.Add(10)
			for j:=0; j<10; j++ {
				go search(i-1)
			}
		} else {
			wg.Add(1)
			go delete(i-1)
		}
		print()
	}

	wg.Wait()
}
