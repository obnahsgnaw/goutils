package errutil

type Builder struct {
	prefix string
}

// NewBuilder
// Deprecated:
func NewBuilder(prefix string) *Builder {
	return &Builder{prefix: prefix}
}

// Deprecated:
func (b *Builder) New(err error, desc ...string) error {
	if err != nil {
		return NewFromError(err, append([]string{b.prefix}, desc...)...)
	}
	return New(append([]string{b.prefix, ": "}, desc...)...)
}

type ErrBuilder struct {
	Prefix string
}

func (b *ErrBuilder) NewError(err error, desc ...string) error {
	if err == nil && len(desc) == 0 {
		return nil
	}
	if b.Prefix == "" {
		panic("error builder prefix not set")
	}
	if err != nil {
		return NewFromError(err, append([]string{b.Prefix}, desc...)...)
	}
	return New(append([]string{b.Prefix, ": "}, desc...)...)
}
