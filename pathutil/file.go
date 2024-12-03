package pathutil

import (
	"github.com/obnahsgnaw/goutils/errutil"
	"os"
	"path/filepath"
)

func ValidFile(filename string) (file string, err error) {
	if filename == "" {
		err = errutil.New("pathutil: filename empty")
		return
	}
	file, err = filepath.Abs(filename)
	if err != nil {
		err = errutil.NewFromError(err, "pathutil: invalid file")
		return
	}
	var stat os.FileInfo
	if stat, err = os.Stat(file); err != nil {
		err = errutil.NewFromError(err, "pathutil: invalid file")
		return
	}
	if stat.IsDir() {
		err = errutil.New("pathutil: not a file=" + filename)
		return
	}

	return
}
