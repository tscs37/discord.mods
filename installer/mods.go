package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

func installMods(basepath string, box *rice.Box) error {
	fmt.Println("Installing embedded mods into discord.mods...")

	var (
		dirCount, filCount int
	)

	err := box.Walk("mods", func(file string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			return os.MkdirAll(filepath.Join(basepath, file), 0755)
		}
		filCount++
		fileS, err := box.Open(file)
		if err != nil {
			return errors.Wrap(err, "Could not open box file")
		}
		defer fileS.Close()
		dst, err := os.OpenFile(filepath.Join(basepath, file), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrap(err, "Could not open target file")
		}
		if err := dst.Truncate(0); err != nil {
			return errors.Wrap(err, "Could not truncate target file")
		}
		defer dst.Close()
		_, err = io.Copy(dst, fileS)
		if err != nil {
			return errors.Wrap(err, "Could not copy file")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "Mod installation failed")
	}

	fmt.Printf("Installed %d files in %d folder.\n", filCount, dirCount)
	return nil
}
