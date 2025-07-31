package example

import (
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"testing"
)

func TestDemo(t *testing.T) {
	d := notifyutil.NewManager()
	d.RegisterChannel(channels.LogChannel, func() notifyutil.Channel {
		return channels.NewLogChannel()
	})
	_, _ = d.DispatchNotification(&DemoNotification{})
}
