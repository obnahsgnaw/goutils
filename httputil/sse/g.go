package sse

import (
	"sync"
)

type Group struct {
	mu      sync.RWMutex
	name    string
	members map[*Client]struct{}
}

func NewGroup(name string) *Group {
	return &Group{
		name:    name,
		members: make(map[*Client]struct{}),
	}
}

func (g *Group) Join(c *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.members[c] = struct{}{}
}

func (g *Group) Leave(c *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.members, c)
}

func (g *Group) Range(handler func(*Client) bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for c := range g.members {
		if !handler(c) {
			return
		}
	}
}

func (g *Group) Members() (list []*Client) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for c := range g.members {
		list = append(list, c)
	}
	return
}

func (g *Group) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.members)
}
