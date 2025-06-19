package arrutil

import "sync"

type StringSet struct {
	sync.Mutex
	sd []string
	md map[string]struct{}
}

func NewStringSet(data []string) *StringSet {
	s := &StringSet{md: make(map[string]struct{})}
	s.Add(data...)
	return s
}

func (s *StringSet) Add(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		if _, ok := s.md[item]; !ok {
			s.md[item] = struct{}{}
			s.sd = append(s.sd, item)
		}
	}
}

func (s *StringSet) Del(items ...string) {
	s.Lock()
	defer s.Unlock()
	hit := false
	for _, item := range items {
		if _, ok := s.md[item]; ok {
			delete(s.md, item)
			hit = true
		}
	}

	if hit {
		s.sd = nil
		for i := range s.md {
			s.sd = append(s.sd, i)
		}
	}
}

func (s *StringSet) Exist(item string) bool {
	_, ok := s.md[item]
	return ok
}

func (s *StringSet) Get() []string {
	return s.sd
}
