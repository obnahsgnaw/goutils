package fileutil

import (
	"bytes"
	"github.com/obnahsgnaw/goutils/arrutil"
	"github.com/obnahsgnaw/goutils/pathutil"
	"os"
	"path/filepath"
	"strings"
)

func ReplaceDir(dir string, fromTo map[string]string, ignore *arrutil.StringSet) (err error) {
	var dirs []os.DirEntry
	if fromTo == nil {
		return
	}
	if dir, err = pathutil.ValidDir(dir); err != nil {
		return
	}
	if dirs, err = os.ReadDir(dir); err != nil {
		return err
	}
	for _, de := range dirs {
		file := filepath.Join(dir, de.Name())
		if strings.HasPrefix(de.Name(), ".") {
			continue
		}
		if de.IsDir() {
			if err = ReplaceDir(file, fromTo, ignore); err != nil {
				return err
			}
		} else {
			if ignore != nil && ignore.Exist(de.Name()) {
				continue
			}
			if err = ReplaceFile(file, fromTo); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReplaceFile(file string, fromTo map[string]string) (err error) {
	var f os.FileInfo
	if fromTo == nil {
		return
	}
	var content []byte
	if f, err = os.Stat(file); err != nil {
		return err
	}
	if content, err = os.ReadFile(file); err != nil {
		return err
	}
	for from, to := range fromTo {
		content = bytes.Replace(content, []byte(from), []byte(to), -1)
	}
	return os.WriteFile(file, content, f.Mode())
}
