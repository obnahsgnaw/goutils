package cmd

import (
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil"
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil/example/cmd"
	"github.com/obnahsgnaw/goutils/singletonutil"
)

func init() {
	StartCommand().RegisterTo(cmd.RootCommand())
}

type Start struct {
	cliutil.BaseCommandBuilder
	config *string
}

func StartCommand() *Start {
	return singletonutil.Instance("cmd:start", func() interface{} {
		c := new(Start)
		c.Initialize(c)
		return c
	}).Get().(*Start)
}

func (s *Start) Name() string {
	return "start"
}

func (s *Start) Desc() string {
	return "start the server"
}

func (s *Start) Handle() {
	println("config path:", *s.config)
	println("Server start running...")
}

func (s *Start) Args(f *cliutil.Flags) {
	s.config = f.StringP("config", "c", "", "config file path")
	f.MarkRequired("config")
}
