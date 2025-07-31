package eventutil

var _dft = NewManger()

// DefaultManger returns a default event manager
func DefaultManger() *Manger {
	return _dft
}
