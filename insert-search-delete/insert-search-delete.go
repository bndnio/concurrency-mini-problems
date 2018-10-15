package main

import (
	"fmt"
	"sync"
	"container/list"
	"./lightswitch"
)

var l *list.List = nil
var wg sync.WaitGroup

insertMutex = &sync.Mutex{}
noSearcher = &sync.Mutex{}
noInserter = &sync.Mutex{}
searchSwitch = &lightswitch
insertSwitch = &lightswitch

func insert(value int) {
	defer wg.Done()
	l.PushBack(value)
}

func search(value int) *list.Element {
	defer wg.Done()
	defer 
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return e
		}
	}
	return nil
}

func delete(value int) {
	defer wg.Done()
	wg.Add(1)
	l.Delete(search(value))
}

func print() {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Print("\n")
}

func main() {
	l = list.New()
	for i:=0; i<20; i++ {
		wg.Add(1)
		insert(i)
	}

	wg.Wait()
}
