package sse

import "context"

var _dft = NewManager(context.Background())

func Default() *Manager {
	return _dft
}

func SetDefaultContent(ctx context.Context) {
	_dft.ctx = ctx
}
