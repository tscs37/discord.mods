package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
	"go.rls.moe/misc/discord.mods/common"
)

var dmods = js.Global.Get("dmodsNS")

func lsDir(path string) ([]string, error) {
	fs, err := common.GetFS()
	if err != nil {
		return nil, err
	}
	filesChan := make(chan []string)
	errChan := make(chan error)

	callback := func(err, files *js.Object) {
		if err.String() != "null" {
			go func() { errChan <- errors.New(err.String()) }()
			go func() { filesChan <- []string{} }()
			return
		}
		go func() { errChan <- nil }()
		var fileList = []string{}
		for i := 0; i < files.Length(); i++ {
			fileList = append(fileList, files.Index(i).String())
		}
		filesChan <- fileList
	}
	go func() {
		fs.Call("readdir", path, callback)
	}()
	return <-filesChan, <-errChan
}

func evalMod(modI *mod) error {
	censored, err := common.CensorPath(modI.baseDir)
	if err != nil {
		return err
	}

	if modI.RunURL != "" {
		print("--> Loading ", modI.Name, " by ", modI.Author, " from ", modI.RunURL)

		for k := range modI.Eval.Pre {
			if err := common.EvalString(modI.Eval.Pre[k]); err != nil {
				return errors.Wrap(err, "Could not eval pre-mod")
			}
		}

		if err := common.EvalURL(modI.RunURL); err != nil {
			return errors.Wrap(err, "Could not eval mod from URL")
		}

	} else {
		print("--> Loading ", modI.Name, " by ", modI.Author, " from ", censored)

		for k := range modI.Eval.Pre {
			if err := common.EvalString(modI.Eval.Pre[k]); err != nil {
				return errors.Wrap(err, "Could not eval pre-mod")
			}
		}
		//indexjsPath := filepath.Join(modI.baseDir, "index.js");
		//mainjsPath := filepath.Join(modI.baseDir, "main.js");
		indexjsPath := common.JoinPath(modI.baseDir, "index.js")
		mainjsPath := common.JoinPath(modI.baseDir, "main.js")
		if exists, err := common.Exists(indexjsPath); exists && err == nil {
			if err := common.RequireFile(modI.baseName, indexjsPath); err != nil {
				return err
			}

		} else if exists, err = common.Exists(mainjsPath); exists && err == nil {
			if err := common.RequireFile(modI.baseName, mainjsPath); err != nil {
				return err
			}
		} else if err != nil {
			return errors.Wrap(err, "Error reading main.js or index.js")
		} else {
			return errors.New("main.js and index.js do not exist")
		}

	}

	for k := range modI.Eval.Post {
		if err := common.EvalString(modI.Eval.Post[k]); err != nil {
			return errors.Wrap(err, "Could not eval post-mod")
		}
	}

	print("--> Finished loading ", modI.Name, " by ", modI.Author, " from ", censored)

	return nil
}
