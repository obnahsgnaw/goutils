package singletonutil

import "sync"

type Builder struct {
	ins       interface{}
	once      sync.Once
	generator func() interface{}
}

func NewBuilder(generator func() interface{}) *Builder {
	return &Builder{generator: generator}
}

func (s *Builder) Get() interface{} {
	s.once.Do(func() {
		s.ins = s.generator()
	})
	return s.ins
}
