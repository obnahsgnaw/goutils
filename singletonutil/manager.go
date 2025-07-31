package singletonutil

import "github.com/obnahsgnaw/goutils/strutil"

type Manager struct {
	prefix   string
	builders map[string]*Builder
}

func NewManager() *Manager {
	return &Manager{builders: make(map[string]*Builder)}
}

func (s *Manager) WithPrefix(prefix string) {
	s.prefix = prefix
}

func (s *Manager) Build(name string, generator func() interface{}) (b *Builder) {
	var ok bool
	name = s.prefixedName(name)
	if b, ok = s.builders[name]; !ok {
		b = NewBuilder(generator)
		s.builders[name] = b
	}

	return b
}

func (s *Manager) prefixedName(name string) string {
	if s.prefix == "" {
		return name
	}
	return strutil.ToString(s.prefix, ":", name)
}
