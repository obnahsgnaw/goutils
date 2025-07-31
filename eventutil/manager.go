package eventutil

import (
	"github.com/asaskevich/EventBus"
)

// Manger event manager
type Manger struct {
	bus EventBus.Bus
}

// NewManger an event manager
func NewManger() *Manger {
	return &Manger{bus: EventBus.New()}
}

// Register an event
func (m *Manger) Register(topic Topic, handler Handler) {
	if topic != "" && handler != nil {
		_ = m.bus.SubscribeAsync(topic.String(), handler, false)
	}
}

// Fire event
func (m *Manger) Fire(e *Event) {
	m.fire(e.Topic, e)
}

func (m *Manger) fire(topic Topic, e *Event) {
	m.bus.Publish(topic.String(), e)
	m.bus.WaitAsync()
}

// Build return an event instance
func (m *Manger) Build(topic Topic, data ...interface{}) *Event {
	return newEvent(topic, data, func() *Manger {
		return m
	})
}
