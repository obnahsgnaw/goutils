package cliutil

import (
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/structutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandBuilder interface {
	Name() string
	Desc() string
	Description() string
	Handle()
	Args(*Flags)
	Build() *cobra.Command
}

type Flags struct {
	cmd *cobra.Command
	*pflag.FlagSet
}

func (f *Flags) MarkRequired(v ...string) {
	for _, i := range v {
		if err := f.cmd.MarkFlagRequired(i); err != nil {
			panic("mark flag required failed: " + err.Error())
		}
	}
}

type BaseCommandBuilder struct {
	structutil.NamedStruct
	errutil.ErrBuilder
	impl CommandBuilder
	cmd  *cobra.Command
}

func (b *BaseCommandBuilder) Initialize(impl CommandBuilder) {
	b.impl = impl
	b.ParseName(b.impl)
	b.ErrPrefix = b.GetName()
}

func (b *BaseCommandBuilder) Build() *cobra.Command {
	if b.impl == nil {
		panic("BaseCliBuilder: cli builder not initialized")
	}
	if b.cmd == nil {
		if b.impl.Name() == "" {
			panic(b.GetName() + ": command name not set")
		}
		b.cmd = &cobra.Command{
			Use:   b.impl.Name(),
			Short: b.impl.Desc(),
			Long:  b.impl.Description(),
			Run: func(_ *cobra.Command, flags []string) {
				b.impl.Handle()
			},
		}
		b.impl.Args(&Flags{b.cmd, b.cmd.Flags()})
	}

	return b.cmd
}

func (b *BaseCommandBuilder) Execute() error {
	return b.Build().Execute()
}

func (b *BaseCommandBuilder) Desc() string {
	return "No Description"
}

func (b *BaseCommandBuilder) Description() string {
	return ""
}

func (b *BaseCommandBuilder) Args(*Flags) {

}

func (b *BaseCommandBuilder) Handle() {

}

func (b *BaseCommandBuilder) RegisterTo(c CommandBuilder) {
	c.Build().AddCommand(b.Build())
}
