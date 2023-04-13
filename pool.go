package main

import (
	"sync"
)

type Pool struct {
	p      chan struct{}
	Wait_g *sync.WaitGroup
	err    error
}

func Newpool(w int) *Pool {
	var e error
	p := make(chan struct{}, w)
	wg := new(sync.WaitGroup)
	return &Pool{p, wg, e}
}

func (p *Pool) worker(f func(arg ...any) error, arg ...any) {
	p.p <- struct{}{}
	p.Wait_g.Add(1)
	go func() {
		defer p.Wait_g.Done()
		err := f(arg...)
		if err != nil {
			p.err = err
		}
		<-p.p
	}()
}
