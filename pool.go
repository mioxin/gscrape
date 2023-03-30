package main

type pool struct {
	p   chan struct{}
	err chan error
}

func Newpool(w int) *pool {
	p := make(chan struct{}, w)
	err := make(chan error)
	return &pool{p, err}
}
func (p *pool) worker(f func(arg ...any) error, arg ...any) {
	p.p <- struct{}{}
	go func() {
		err := f(arg...)
		if err != nil {
			p.err <- err
		}
		<-p.p
	}()
}
