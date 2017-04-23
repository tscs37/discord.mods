package osmode

import (
	"io/ioutil"
	"os"
)

func GetFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func WriteFile(file string, content []byte, mode int) error {
	return ioutil.WriteFile(file, content, os.FileMode(mode))
}

func Mkdir(path string, mode int) error {
	return os.MkdirAll(path, os.FileMode(mode))
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if os.IsExist(err) {
		return true, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}
