package main

import (
	"fmt"
	"sync"
)

var (
	p    *people
	once sync.Once
)

func init() {
	once.Do(func() {
		p = &people{}
	})
}

type people struct {
	age int
	mux sync.Mutex
}

func GetPeopleInstance() *people {
	return p
}

//这里的更改操作必须加锁保证线程安全
func (p *people) IncreaseAge() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.age++
}
func (p *people) GetAge() int {
	return p.age
}
func main() {
	wg := sync.WaitGroup{}
	wg.Add(30000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			p1 := GetPeopleInstance()
			p1.IncreaseAge()
		}()
		go func() {
			defer wg.Done()
			p2 := GetPeopleInstance()
			p2.IncreaseAge()
		}()
		go func() {
			defer wg.Done()
			p3 := GetPeopleInstance()
			p3.IncreaseAge()
		}()
	}
	wg.Wait()
	p4 := GetPeopleInstance()
	fmt.Println(p4.age)
}
