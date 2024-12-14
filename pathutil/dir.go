package pathutil

import (
	"github.com/obnahsgnaw/goutils/errutil"
	"os"
	"path/filepath"
)

func ValidDir(dirname string) (dir string, err error) {
	if dirname == "" {
		err = errutil.New("pathutil: dirname empty")
		return
	}
	dir, err = filepath.Abs(dirname)
	if err != nil {
		err = errutil.NewFromError(err, "pathutil: invalid path")
		return
	}
	var stat os.FileInfo
	if stat, err = os.Stat(dir); err != nil {
		err = errutil.NewFromError(err, "pathutil: invalid path")
		return
	}
	if !stat.IsDir() {
		err = errutil.New("pathutil: not a directory=", dirname)
		return
	}

	return
}

func RuntimeDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errutil.NewFromError(err, "pathutil: cannot get executable path")
	}
	dir := filepath.Dir(ex)
	return dir, nil
}

func Mkdir(dirPath string, mode os.FileMode) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.Mkdir(dirPath, mode); err != nil {
			return err
		}
	}
	return nil
}

func MkdirAll(dirPath string, mode os.FileMode) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, mode); err != nil {
			return err
		}
	}
	return nil
}
