package pathutil

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func CopyEmbedFsDir(f embed.FS, dirName, rootPath string, replace func(name string, content []byte) (string, []byte)) (err error) {
	if replace == nil {
		replace = func(name string, content []byte) (string, []byte) { return name, content }
	}
	rpDir, _ := replace(dirName, nil)
	if rpDir == "" {
		return nil
	}
	if err = MkdirAll(path.Join(rootPath, rpDir), 0755); err != nil {
		return
	}
	var dirs []fs.DirEntry
	if dirs, err = f.ReadDir(dirName); err != nil {
		return
	}
	for _, d := range dirs {
		if d.IsDir() {
			if err = CopyEmbedFsDir(f, path.Join(dirName, d.Name()), rootPath, replace); err != nil {
				return
			}
		} else {
			fileName := path.Join(dirName, d.Name())
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
	if err = os.WriteFile(filepath.Join(rootPath, rpFileName), rpContent, 0640); err != nil {
		return
	}
	return
}
