package notifyutil

import "sync"

type channelBuilder struct {
	handler func() Channel
	once    sync.Once
	ins     Channel
}

func (b *channelBuilder) Get() Channel {
	b.once.Do(func() {
		b.ins = b.handler()
	})
	return b.ins
}

type Manager struct {
	channels map[ChannelName]*channelBuilder
}

func NewManager() *Manager {
	return &Manager{
		channels: make(map[ChannelName]*channelBuilder),
	}
}

// RegisterChannel register a channel
func (d *Manager) RegisterChannel(channel ChannelName, handler func() Channel) {
	if handler != nil {
		d.channels[channel] = &channelBuilder{
			handler: handler,
		}
	}
}

// DispatchNotification Distribution Notices
// dispatched: If there is a successful sending, it is true
// lastErr: The last error
func (d *Manager) DispatchNotification(notification Notification) (dispatched bool, lastErr error) {
	for _, ch := range notification.Channels() {
		if handler, ok := d.channels[ch]; ok {
			if notification.RetryMax() <= 0 {
				continue
			}
			for retry := 1; retry <= notification.RetryMax(); retry++ {
				chData := notification.Data(ch)
				if notification.To(ch) == "" || chData == nil {
					break
				}
				if lastErr = handler.Get().Send(notification.To(ch), chData); lastErr == nil {
					break
				}
			}
			if lastErr != nil {
				notification.FailedHandle(ch, lastErr)
			} else {
				dispatched = true
				notification.SuccessHandle(ch)
			}
		}
	}
	return
}
