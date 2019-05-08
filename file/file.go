package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func GetAllFiles(dirPth string) ([]string, error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, fi := range dir {
		if !fi.IsDir() {
			if filepath.Ext(fi.Name()) == ".cert" {
				files = append(files, filepath.Join(dirPth, fi.Name()))
			}
		}
	}
	return files, nil
}
