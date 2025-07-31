package sse

import (
	"context"
	"sync"
)

const (
	Connect Event = iota
	Disconnect
)

// Manager Maintain SSE connections
type Manager struct {
	ctx      context.Context
	mu       sync.Mutex
	clients  map[*Client]struct{}
	tags     map[string]*Group
	groups   map[string]*Group
	listener func(e Event, c *Client)
}

type Event int

func NewManager() *Manager {
	return &Manager{
		ctx:     context.Background(),
		clients: make(map[*Client]struct{}),
		tags:    make(map[string]*Group),
		groups:  make(map[string]*Group),
		listener: func(e Event, c *Client) {

		},
	}
}

func (s *Manager) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Manager) NewClient() *Client {
	c := NewClient(s.ctx)
	c.service = s
	s.addClient(c)
	return c
}

func (s *Manager) AddClient(c *Client) {
	c.service = s
	s.addClient(c)
}

func (s *Manager) addClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = struct{}{}
	s.listener(Connect, c)
}

func (s *Manager) GetClient(tag string) []*Client {
	v, ok := s.tags[tag]
	if ok {
		return v.Members()
	}
	return nil
}

func (s *Manager) RemoveClient(c *Client) {
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
	s.listener(Disconnect, c)
}

func (s *Manager) BroadcastAll(message *Message) {
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

func (s *Manager) BroadcastGroup(group string, message *Message) {
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

func (s *Manager) BroadcastTag(tag string, message *Message) {
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

func (s *Manager) ConnectionCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}

func (s *Manager) GroupCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.groups)
}

func (s *Manager) SetListener(listener func(e Event, c *Client)) {
	if listener != nil {
		s.listener = listener
	}
}

func (s *Manager) addClientTag(c *Client, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tags[tag]; !ok {
		s.tags[tag] = NewGroup(tag)
	}
	s.tags[tag].Join(c)
}

func (s *Manager) rmClientTag(c *Client, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.tags[tag]; ok {
		v.Leave(c)
		if v.Len() == 0 {
			delete(s.tags, tag)
		}
	}
}

func (s *Manager) addClientGroup(c *Client, group string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[group]; !ok {
		s.groups[group] = NewGroup(group)
	}
	s.groups[group].Join(c)
}

func (s *Manager) rmClientGroup(c *Client, group string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.groups[group]; ok {
		v.Leave(c)
		if v.Len() == 0 {
			delete(s.groups, group)
		}
	}
}
