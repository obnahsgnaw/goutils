package notifyutil

type Notification interface {
	Channels() []ChannelName
	To(ChannelName) Target
	Data(ChannelName) Data
	SuccessHandle(ChannelName)
	FailedHandle(ChannelName, error)
	RetryMax() int
}
