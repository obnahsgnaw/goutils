package pathutil

import (
	"embed"
	"io/fs"
	"os"
	"strings"
)

func CopyEmbedFsDir(f embed.FS, dirName, rootPath string, replace func(name string, content []byte) (string, []byte)) (err error) {
	if replace == nil {
		replace = func(name string, content []byte) (string, []byte) { return name, content }
	}
	rpDir, _ := replace(dirName, nil)
	if rpDir == "" {
		return nil
	}
	if err = MkdirAll(filepathJoin(rootPath, rpDir), 0755); err != nil {
		return
	}
	var dirs []fs.DirEntry
	if dirs, err = f.ReadDir(dirName); err != nil {
		return
	}
	for _, d := range dirs {
		if d.IsDir() {
			if err = CopyEmbedFsDir(f, filepathJoin(dirName, d.Name()), rootPath, replace); err != nil {
				return
			}
		} else {
			fileName := filepathJoin(dirName, d.Name())
			if err = CopyEmbedFsFile(f, fileName, rootPath, replace); err != nil {
				return
			}
		}
	}
	return
}

func CopyEmbedFsFile(f embed.FS, fileName, rootPath string, replace func(name string, content []byte) (string, []byte)) (err error) {
	if replace == nil {
		replace = func(name string, content []byte) (string, []byte) { return name, content }
	}
	var content []byte
	if content, err = f.ReadFile(fileName); err != nil {
		return
	}
	rpFileName, rpContent := replace(fileName, content)
	if rpFileName == "" {
		return
	}
	if err = os.WriteFile(filepathJoin(rootPath, rpFileName), rpContent, 0640); err != nil {
		return
	}
	return
}

func filepathJoin(name ...string) string {
	return strings.Join(name, "/")
}
