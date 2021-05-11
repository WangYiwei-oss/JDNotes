package main

import (
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
)

type SpinMutex int32

func (s *SpinMutex) Lock() {
	for !atomic.CompareAndSwapInt32((*int32)(s), 0, 1) {
	}
}
func (s *SpinMutex) UnLock() {
	atomic.CompareAndSwapInt32((*int32)(s), 1, 0)
}

type Locker interface {
	Lock()
	UnLock()
}

var wg = sync.WaitGroup{}

func routine(flag int, n *int, l Locker, time int) {
	for i := 0; i < time; i++ {
		l.Lock()
		*n++
		l.UnLock()
	}
	wg.Done()
}
func main() {
	var p bool
	flag.BoolVar(&p, "pause", false, "with pause instruction")
	flag.Parse()
	l := new(SpinMutex)
	n := 0
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go routine(i, &n, l, 100000)
	}
	wg.Wait()
	fmt.Println(n)
}
