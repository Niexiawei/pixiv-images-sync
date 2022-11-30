package gpool

import "sync"

type GPool struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

func NewGPool(maxSize int) *GPool {
	return &GPool{
		c:  make(chan struct{}, maxSize),
		wg: new(sync.WaitGroup),
	}
}

func (s *GPool) Add(delta int) {
	s.wg.Add(delta)
	if delta > 1 {
		for i := 0; i < delta; i++ {
			s.c <- struct{}{}
		}
	} else {
		s.c <- struct{}{}
	}
}

func (s *GPool) Done() {
	s.wg.Done()
	<-s.c
}

func (s *GPool) Wait() {
	s.wg.Wait()
}
