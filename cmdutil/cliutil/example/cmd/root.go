package cmd

import (
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil"
	"github.com/obnahsgnaw/goutils/singletonutil"
)

type Root struct {
	cliutil.BaseCommandBuilder
}

func RootCommand() *Root {
	return singletonutil.Instance("cmd:root", func() interface{} {
		c := new(Root)
		c.Initialize(c)
		return c
	}).Get().(*Root)
}

func (s *Root) Name() string {
	return "demo"
}

func (s *Root) Description() string {
	return `this is command description`
}
