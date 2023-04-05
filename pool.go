package main

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	p      chan struct{}
	err    chan error
	Wait_g *sync.WaitGroup
}

func Newpool(w int) *Pool {
	p := make(chan struct{}, w)
	err := make(chan error)
	wg := new(sync.WaitGroup)
	return &Pool{p, err, wg}
}

func (p *Pool) worker(f func(arg ...any) error, arg ...any) {
	p.p <- struct{}{}
	p.Wait_g.Add(1)
	fmt.Println("worker get chan", p.Wait_g)
	go func() {
		f(arg...)
		// if err != nil {
		// 	p.err <- err
		// }
		<-p.p
		p.Wait_g.Done()
		fmt.Println("goroutin out chan")
	}()
	time.Sleep(30 * time.Second)
}
