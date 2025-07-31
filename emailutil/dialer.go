package emailutil

type Dialer interface {
	DialAndSend(e ...*Email) error
	Ssl(ssl bool)
}
