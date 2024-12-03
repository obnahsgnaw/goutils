package main

import "github.com/obnahsgnaw/goutils/pathutil"

func main() {
	currentDir, err := pathutil.RuntimeDir()
	if err != nil {
		panic(err)
	}
	println(currentDir)

	dir, err := pathutil.ValidDir("./")
	if err != nil {
		panic(err)
	}
	println(dir)

	f, err1 := pathutil.ValidFile("./path.go")
	if err1 != nil {
		panic(err1)
	}
	println(f)
}
