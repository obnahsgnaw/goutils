package errutil

type Builder struct {
	prefix string
}

func NewBuilder(prefix string) *Builder {
	return &Builder{prefix: prefix}
}

func (b *Builder) New(err error, desc ...string) error {
	if err != nil {
		return NewFromError(err, append([]string{b.prefix}, desc...)...)
	}
	return New(append([]string{b.prefix, ": "}, desc...)...)
}
