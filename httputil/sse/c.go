package sse

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Client Identify an SSE client
type Client struct {
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	service  *Manager
	w        http.ResponseWriter
	r        *http.Request
	optional interface{}
	channel  chan *Message
	tags     map[string]struct{}
	groups   map[string]struct{}
}

// NewClient Initialize a new connection, which can also be maintained in the Manager
func NewClient(ctx context.Context) *Client {
	c1, cl := context.WithCancel(ctx)
	return &Client{
		ctx:     c1,
		cancel:  cl,
		channel: make(chan *Message, 5),
		tags:    make(map[string]struct{}),
		groups:  make(map[string]struct{}),
	}
}

// SetOptional Set up a custom binding data object
func (c *Client) SetOptional(optional interface{}) {
	c.optional = optional
}

// GetOptional Returns the custom data binding object that was set
func (c *Client) GetOptional() interface{} {
	return c.optional
}

// AddTag Add custom binding tags, and later sections query based on the tags in the Manager
func (c *Client) AddTag(tag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tags[tag] = struct{}{}
	if c.service != nil {
		c.service.addClientTag(c, tag)
	}
}

// RmTag Remove custom tags
func (c *Client) RmTag(tag string) {
	delete(c.tags, tag)
	if c.service != nil {
		c.service.rmClientTag(c, tag)
	}
}

// JoinGroup Join the group, and the follow-up department will broadcast within the group
func (c *Client) JoinGroup(group string) {
	c.groups[group] = struct{}{}
	if c.service != nil {
		c.service.addClientGroup(c, group)
	}
}

// LeaveGroup Leave the group
func (c *Client) LeaveGroup(group string) {
	delete(c.groups, group)
	if c.service != nil {
		c.service.rmClientGroup(c, group)
	}
}

// Start HTTP service processing
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if c.service != nil {
			c.service.RemoveClient(c)
		}
		close(c.channel)
	}()
	c.w = w
	c.r = r
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	go func() {
		for {
			select {
			case <-r.Context().Done():
				c.cancel()
				return
			case <-c.ctx.Done():
				return
			case msg := <-c.channel:
				c.write(msg)
			case <-time.After(time.Second * 10):
			}
		}
	}()
	<-c.ctx.Done()
}

// Send Send data messages
func (c *Client) Send(message *Message) {
	c.channel <- message
}

func (c *Client) write(message *Message) {
	if c.w == nil {
		return
	}

	_, err := c.w.Write([]byte(Encode(message)))
	if err != nil {
		c.cancel()
	}
	c.w.(http.Flusher).Flush()
}

// Provider Set the timing of the data supply
func (c *Client) Provider(provider func() *Message, interval time.Duration) {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if msg := provider(); msg != nil {
					c.Send(msg)
				}
				time.Sleep(interval)
			}
		}
	}()
}
