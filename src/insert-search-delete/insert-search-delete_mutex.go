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
	"lightswitch_mutex"
)

var l *list.List = nil
var wg sync.WaitGroup

var insertMutex = &sync.Mutex{}
var noSearcher = &sync.Mutex{}
var noInserter = &sync.Mutex{}
var searchSwitch = lightswitch_mutex.New()
var insertSwitch = lightswitch_mutex.New()

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
	insertMutex.Lock()
	found := find(value)
	insertMutex.Unlock()
	insertSwitch.Unlock(noInserter)

	return found
}

func delete(value int) {
	defer wg.Done()
	noSearcher.Lock()
	noInserter.Lock()
	l.Remove(find(value))
	noInserter.Unlock()
	noSearcher.Unlock()
}

func print() {
	for e := l.Front(); e != nil; e = e.Next() {
		// fmt.Print(e.Value, " ")
	}
	// fmt.Print("\n")
}

func main() {

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
