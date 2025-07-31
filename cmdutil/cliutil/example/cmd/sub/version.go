package cmd

import (
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil"
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil/example/cmd"
	"github.com/obnahsgnaw/goutils/singletonutil"
)

func init() {
	VersionCommand().RegisterTo(cmd.RootCommand())
}

type Version struct {
	cliutil.BaseCommandBuilder
}

func VersionCommand() *Version {
	return singletonutil.Instance("cmd:sub:version", func() interface{} {
		c := new(Version)
		c.Initialize(c)
		return c
	}).Get().(*Version)
}

func (s *Version) Name() string {
	return "version"
}

func (s *Version) Desc() string {
	return "show the server version"
}

func (s *Version) Handle() {
	println("Version: v1.0.0")
}
