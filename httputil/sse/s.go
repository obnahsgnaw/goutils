package sse

import (
	"context"
	"sync"
)

type Service struct {
	ctx      context.Context
	mu       sync.Mutex
	clients  map[*Client]struct{}
	tags     map[string]*Group
	groups   map[string]*Group
	listener func(e Event, c *Client)
}

type Event int

const (
	Connect Event = iota
	DisConnect
)

func New(ctx context.Context) *Service {
	return &Service{
		ctx:     ctx,
		clients: make(map[*Client]struct{}),
		tags:    make(map[string]*Group),
		groups:  make(map[string]*Group),
		listener: func(e Event, c *Client) {

		},
	}
}

func (s *Service) NewClient() *Client {
	c := NewClient(s.ctx)
	c.service = s
	s.addClient(c)
	return c
}

func (s *Service) AddClient(c *Client) {
	c.service = s
	s.addClient(c)
}

func (s *Service) GetClient(tag string) []*Client {
	v, ok := s.tags[tag]
	if ok {
		return v.Members()
	}
	return nil
}

func (s *Service) RemoveClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c)
	for tag := range c.tags {
		if v, ok := s.tags[tag]; ok {
			v.Leave(c)
			if v.Len() == 0 {
				delete(s.tags, tag)
			}
		}
	}
	for group := range c.groups {
		if v, ok := s.groups[group]; ok {
			v.Leave(c)
			if v.Len() == 0 {
				delete(s.groups, group)
			}
		}
	}
	s.listener(DisConnect, c)
}

func (s *Service) BroadcastAll(message *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var wg sync.WaitGroup
	if len(s.clients) == 0 {
		return
	}
	for c := range s.clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			c.Send(message)
		}(c)
	}
	wg.Wait()
}

func (s *Service) BroadcastGroup(group string, message *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var wg sync.WaitGroup
	if len(s.clients) == 0 {
		return
	}
	if group == "" {
		return
	}
	if v, ok := s.groups[group]; ok {
		v.Range(func(cc *Client) bool {
			wg.Add(1)
			go func(c *Client) {
				defer wg.Done()
				c.Send(message)
			}(cc)
			return true
		})
	}
	wg.Wait()
}

func (s *Service) BroadcastTag(tag string, message *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var wg sync.WaitGroup
	if len(s.clients) == 0 {
		return
	}
	if tag == "" {
		return
	}
	if v, ok := s.tags[tag]; ok {
		v.Range(func(cc *Client) bool {
			wg.Add(1)
			go func(c *Client) {
				defer wg.Done()
				c.Send(message)
			}(cc)
			return true
		})
	}
	wg.Wait()
}

func (s *Service) ConnectionCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}

func (s *Service) GroupCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.groups)
}

func (s *Service) SetListener(listener func(e Event, c *Client)) {
	if listener != nil {
		s.listener = listener
	}
}

func (s *Service) addClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = struct{}{}
	s.listener(Connect, c)
}

func (s *Service) addClientTag(c *Client, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tags[tag]; !ok {
		s.tags[tag] = NewGroup(tag)
	}
	s.groups[tag].Join(c)
}

func (s *Service) rmClientTag(c *Client, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.tags[tag]; ok {
		v.Leave(c)
		if v.Len() == 0 {
			delete(s.tags, tag)
		}
	}
}

func (s *Service) addClientGroup(c *Client, group string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[group]; !ok {
		s.groups[group] = NewGroup(group)
	}
	s.groups[group].Join(c)
}

func (s *Service) rmClientGroup(c *Client, group string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.groups[group]; ok {
		v.Leave(c)
		if v.Len() == 0 {
			delete(s.groups, group)
		}
	}
}
