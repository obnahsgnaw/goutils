package example

import (
	"github.com/obnahsgnaw/goutils/emailutil"
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"testing"
)

func TestDemo1(t *testing.T) {
	notifyutil.DefaultManager().RegisterChannel(channels.LogChannel, func() notifyutil.Channel {
		return channels.NewLogChannel()
	})
	notifyutil.DefaultManager().RegisterChannel(channels.EmailChannel, func() notifyutil.Channel {
		return channels.NewEmailChannel(emailutil.NewDevManager())
	})
	n := NewRegisteredNotification("Jo", "xxx@xx.com")
	n.RegisterTo(notifyutil.DefaultManager())
	n.Notify()
}
