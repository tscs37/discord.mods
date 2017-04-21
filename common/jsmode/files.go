package jsmode

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
	"go.rls.moe/misc/discord.mods/common"
)

func GetFile(filename string) ([]byte, error) {
	fs, err := GetModule("fs")
	if err != nil {
		return nil, errors.Wrap(err, "Could not acquire FS module")
	}
	if exists, err := common.Exists(filename); !exists || err != nil {
		if err != nil {
			return nil, errors.Wrap(err, "Error while checking file")
		}
		return nil, errors.New("Does not exist")
	}
	datChan := make(chan string)
	errChan := make(chan error)

	callback := func(err, data *js.Object) {
		if err.String() == "null" {
			go func() { errChan <- nil }()
		} else {
			go func() { errChan <- errors.New(err.String()) }()
		}
		datChan <- data.String()
	}

	go func() {
		fs.Call("readFile", filename, callback)
	}()

	return []byte(<-datChan), <-errChan
}

func WriteFile(file string, content []byte, mode int) error {
	errChan := make(chan error)
	callback := func(err *js.Object) {
		if err.String() == "null" {
			errChan <- nil
			return
		}
		errChan <- errors.New(err.String())
	}
	go func() {
		fs, err := common.GetFS()
		if err != nil {
			errChan <- err
			return
		}
		fs.Call("writeFile", file, content, map[string]interface{}{
			"mode": mode,
		}, callback)
	}()
	return <-errChan
}

func Mkdir(path string, mode int) error {
	fs, err := common.GetFS()
	if err != nil {
		return err
	}
	fs.Call("mkdirSync", path, mode)
	return nil
}

func Exists(path string) (bool, error) {
	fs, err := common.GetFS()
	if err != nil {
		return false, err
	}
	return fs.Call("existsSync", path).Bool(), nil
}
