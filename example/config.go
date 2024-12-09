package main

import (
	"fmt"
	"github.com/obnahsgnaw/goutils/configutil"
	"os"
)

type MyConfig struct {
	configutil.BaseConfig
	Version         bool `short:"v" long:"version" description:"show version"`
	versionProvider func() string
	ConfigFile      string `short:"c" long:"conf" description:"config file"`

	A string `short:"a" long:"aa" description:"show aa"`
}

func (s *MyConfig) Parse() error {
	if err := s.ParseFlag(s); err != nil {
		return err
	}

	if s.Version {
		if s.versionProvider != nil {
			fmt.Println(s.versionProvider())
		} else {
			fmt.Println("Unknown, version provider not set")
		}
		os.Exit(0)
	}

	return s.ParseFile(s, s.ConfigFile)
}

func main() {
	cnf := MyConfig{}
	if err := cnf.Parse(); err != nil {
		panic(err)
	}
	println(cnf.A)
}
