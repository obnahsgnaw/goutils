package eventutil

import (
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/structutil"
)

type EventBuilder interface {
	Topic() Topic
	Handle(*Event)
}

type BaseEventBuilder struct {
	structutil.NamedStruct
	errutil.ErrBuilder
	impl    EventBuilder
	manager *Manger
}

func (s *BaseEventBuilder) Fire(data ...interface{}) {
	if s.impl == nil {
		panic("BaseEventBuilder: Initialize() first when event Fire()")
	}
	if s.manager == nil {
		panic("BaseEventBuilder: RegisterTo() manager first when event Fire()")
	}
	s.manager.Build(s.impl.Topic(), data...).Fire()
}

func (s *BaseEventBuilder) Initialize(impl EventBuilder) {
	s.impl = impl
	s.ParseName(s.impl)
	s.ErrPrefix = s.GetName()
	if s.manager != nil {
		s.manager.Register(s.impl.Topic(), s.impl.Handle)
	}
}

func (s *BaseEventBuilder) RegisterTo(m *Manger) {
	s.manager = m
	if s.impl != nil {
		s.manager.Register(s.impl.Topic(), s.impl.Handle)
	}
}
