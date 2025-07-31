package notifyutil

import (
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/structutil"
)

type BaseNotificationBuilder struct {
	structutil.NamedStruct
	errutil.ErrBuilder
	impl    Notification
	manager *Manager
}

func (s *BaseNotificationBuilder) Initialize(impl Notification) {
	s.impl = impl
	s.ParseName(s.impl)
	s.ErrPrefix = s.GetName()
}

func (s *BaseNotificationBuilder) RegisterTo(m *Manager) {
	s.manager = m
}

func (s *BaseNotificationBuilder) Notify() {
	if s.impl == nil {
		panic("BaseNotificationBuilder: Initialize() first when notification Notify()")
	}
	if s.manager == nil {
		panic("BaseNotificationBuilder: RegisterTo() manager first when notification Notify()")
	}
	_, _ = s.manager.DispatchNotification(s.impl)
	return
}
