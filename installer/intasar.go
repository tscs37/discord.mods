package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"layeh.com/asar"
	"os"
	"strings"
)

func runIntAsar(path, bootstrapper string) error {
	fmt.Println(`

	You are using the internal ASAR engine,
	be aware that this engine is experimental
	and may brick your discord install.

	`)

	file, err := os.Open(getResourcePath(path))
	if err != nil {
		return errors.Wrap(err, "Could not open ASAR")
	}
	defer file.Close()

	discordAsar, err := asar.Decode(file)
	if err != nil {
		return errors.Wrap(err, "Could not decode ASAR")
	}

	fmt.Println("Beginning install...")
	for k := range discordAsar.Children {
		file := discordAsar.Children[k]
		if file.Path() == "index.js" {
			oldIndexJsStr := file.String()
			oldIndexJsStr, err = uninstall(oldIndexJsStr)
			if err != nil {
				return errors.Wrap(err, "Could not uninstall bootstrap.js")
			}
			newIndexJs, err := install(oldIndexJsStr, bootstrapper)
			if err != nil {
				return errors.Wrap(err, "Could not install bootstrap.js")
			}
			newEntry := asar.New("index.js", bytes.NewReader([]byte(newIndexJs)), int64(len(newIndexJs)), 0, file.Flags)
			discordAsar.Children[k] = newEntry
			break
		}
	}

	outfile, err := os.Create("./output.asar")
	if err != nil {
		return errors.Wrap(err, "Output file could not be created")
	}
	defer outfile.Close()

	builder := &asar.Builder{}

	err = discordAsar.Walk(func(file string, info os.FileInfo, err error) error {
		recBuilder(discordAsar, builder, file, info, nil)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Could not fully walk ASAR Contents")
	}

	_, err = builder.Root().EncodeTo(outfile)
	if err != nil {
		return errors.Wrap(err, "Could not encode ASAR")
	}
	return nil
}

func recBuilder(entry *asar.Entry, builder *asar.Builder, file string, info os.FileInfo, content *string) {
	fileSegs := strings.Split(file, "/")
	recBuilderR(entry, builder, fileSegs, info, content)
}

func recBuilderR(entry *asar.Entry, builder *asar.Builder, file []string, info os.FileInfo, content *string) {
	if len(file) == 0 {
		return
	}
	if len(file) == 1 {
		if info.IsDir() {
			builder.AddDir(file[0], asar.Flag(info.Mode()))
			return
		} else {
			var cBytes = []byte{}
			if content != nil {
				cBytes = bytes.NewBufferString(*content).Bytes()
			} else {
				cBytes = entry.Find(file[0]).Bytes()
			}
			builder.Add(file[0], bytes.NewReader(cBytes), int64(len(cBytes)), asar.Flag(info.Mode()))
			return
		}
	}
	newBuilder := builder.AddDir(file[0], entry.Find(file[0]).Flags)
	recBuilderR(entry.Find(file[0]), newBuilder, file[1:], info, content)
}
