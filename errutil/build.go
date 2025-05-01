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
	ErrPrefix string
}

func (b *ErrBuilder) NewError(err error, desc ...string) error {
	if err == nil && len(desc) == 0 {
		return nil
	}
	if b.ErrPrefix == "" {
		panic("error builder prefix not set")
	}
	prefixedDesc := []string{b.ErrPrefix}
	if len(desc) > 0 {
		prefixedDesc = append(prefixedDesc, ": ")
		prefixedDesc = append(prefixedDesc, desc...)
	}
	if err != nil {
		return NewFromError(err, prefixedDesc...)
	}
	return New(prefixedDesc...)
}
