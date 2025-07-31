package singletonutil

var _defaultManager *Manager

func init() {
	_defaultManager = NewManager()
}

func Default() *Manager {
	return _defaultManager
}

func Instance(name string, gen func() interface{}) *Builder {
	return _defaultManager.Build(name, gen)
}
