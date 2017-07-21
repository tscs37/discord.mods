package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func runExtAsar(path, bootstrapper string) error {
	fmt.Println(`

	You are using the external ASAR engine,
	please ensure you have a compatible asar binary installed

	`)

	tmp, err := unpackAsar(filepath.Join(getResourcePath(path), "app.asar"))
	if err != nil {
		return errors.Wrap(err, "Unpacking failed")
	}

	var oldFlags os.FileMode
	if stat, err := os.Stat(filepath.Join(tmp, "index.js")); os.IsNotExist(err) {
		return errors.Wrap(err, "Could not find index.js in unpacked ASAR")
	} else {
		oldFlags = stat.Mode()
	}

	oldIndexJs, err := ioutil.ReadFile(filepath.Join(tmp, "index.js"))
	if err != nil {
		return errors.Wrap(err, "Could not read index.js")
	}
	oldIndexJsStr := string(oldIndexJs)

	fmt.Println("Beginning install...")
	oldIndexJsStr, err = uninstall(oldIndexJsStr)
	if err != nil {
		return errors.Wrap(err, "Could not uninstall bootstrap.js")
	}
	newIndexJs, err := install(oldIndexJsStr, bootstrapper)
	if err != nil {
		return errors.Wrap(err, "Could not install bootstrap.js")
	}

	err = ioutil.WriteFile(filepath.Join(tmp, "index.js"), []byte(newIndexJs), oldFlags)
	if err != nil {
		return errors.Wrap(err, "Could not writeback index.js")
	}

	err = packAsar(tmp, "./output.asar")
	if err != nil {
		return errors.Wrap(err, "Could not pack ASAR")
	}
	return nil
}

func unpackAsar(asar string) (string, error) {
	tmp, err := ioutil.TempDir("", "discord.mods")
	if err != nil {
		return "", errors.Wrap(err, "Error on TempDir creation")
	}
	cmd := exec.Command("asar", "extract", asar, tmp)
	err = cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, "ASAR did not extract")
	}
	return tmp, nil
}

func packAsar(dir, asar string) error {
	cmd := exec.Command("asar", "pack", dir, asar)
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "ASAR did not pack")
	}
	return nil
}
