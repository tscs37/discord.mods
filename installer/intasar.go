package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"layeh.com/asar"
	"os"
	"path/filepath"
	"strings"
)

func runIntAsar(path, bootstrapper string) error {
	fmt.Println(`

	You are using the internal ASAR engine,
	be aware that this engine is experimental
	and may brick your discord install.

	`)

	if err := unpackAsarInt(path, getResourcePath(path)); err != nil {
		return errors.Wrap(err, "Could not unpack ASAR")
	}

	indexJsPath := filepath.Join(getResourcePath(path), "index.js")

	stat, err := os.Stat(indexJsPath)
	if err != nil {
		return errors.Wrap(err, "COuld not stat index.js")
	}
	data, err := ioutil.ReadFile(indexJsPath)
	if err != nil {
		return errors.Wrap(err, "Could not read index.js")
	}

	newIndexJs, err := install(string(data), bootstrapper)
	if err != nil {
		return errors.Wrap(err, "Could not patch index.js")
	}

	if err := ioutil.WriteFile(indexJsPath, []byte(newIndexJs), stat.Mode()); err != nil {
		return errors.Wrap(err, "Could not write index.js")
	}

	return nil
}

func unpackAsarInt(asarPath, target string) error {
	file, err := os.Open(filepath.Join(getResourcePath(asarPath), "app.asar"))
	if err != nil {
		return errors.Wrap(err, "Could not open ASAR")
	}
	asarFile, err := asar.Decode(file)
	if err != nil {
		return errors.Wrap(err, "Could not decode ASAR")
	}
	return asarFile.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(target, path), info.Mode())
		}
		ent := asarFile.Find(strings.Split(path, "/")...)
		if ent == nil {
			return errors.Errorf("Could not find file %s in ASAR", path)
		}
		return ioutil.WriteFile(filepath.Join(target, path), ent.Bytes(), info.Mode())
	})
}
