package main

import (
	"github.com/obnahsgnaw/goutils/cmdutil/cliutil/example/cmd"
	_ "github.com/obnahsgnaw/goutils/cmdutil/cliutil/example/cmd/sub"
)

func main() {
	_ = cmd.RootCommand().Execute()
}
