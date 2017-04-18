package main

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

var dmods = js.Global.Get("dmodsNS")

func evalFile(file string) error {
	dat, err := getFile(file)
	if err != nil {
		return err
	}
	evalString(dat)
	return nil
}

func evalString(stmt string) {
	js.Global.Call("eval", stmt)
}

func requireFile(varname, file string) {
	dmods.Set(varname, js.Global.Call("require", file))
}

func evalUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		print("Error on loading resource: ", err)
		return
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print("Error on reading resource: ", err)
		return
	}
	evalString(string(dat))
}

func getJqueryReliable() string {
	if js.Global.Get("jQuery").String() == "undefined" {
		return "undefined"
	}
	if js.Global.Get("jQuery").Get("fn").String() == "undefined" {
		return "undefined"
	}
	return js.Global.Get("jQuery").Get("fn").Get("jquery").String()
}

func getFile(file string) (string, error) {
	fs := getFs()
	if !fileExists(file) {
		return "", errors.New("Does not exist")
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
		fs.Call("readFile", file, callback)
	}()

	return <-datChan, <-errChan
}

func fileExists(file string) bool {
	return getFs().Call("existsSync", file).Bool()
}

func writeFile(file, content string, mode int) error {
	errChan := make(chan error)
	callback := func(err *js.Object) {
		if err.String() == "null" {
			errChan <- nil
			return
		}
		errChan <- errors.New(err.String())
	}
	go func() {
		getFs().Call("writeFile", file, content, map[string]interface{}{
			"mode": mode,
		}, callback)
	}()
	return <-errChan
}

func mkdirPath(path string, mode int) {
	getFs().Call("mkdirSync", path, mode)
}

func dirExists(path string) bool {
	return getFs().Call("existsSync", path).Bool()
}

func lsDir(path string) ([]string, error) {
	fs := getFs()
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

func getHome() string {
	return dmods.Call("homedir").String()
}

func getBase() string {
	return filepath.Join(getHome(), baseFolder)
}

func alert(message string) {
	js.Global.Call("alert", message)
}

func getFs() *js.Object {
	return dmods.Call("getFS")
}

func censorPath(path string) string {
	return strings.TrimPrefix(path, getHome())
}

func evalMod(modI *mod) error {
	if modI.RunURL != "" {
		print("--> Loading ", modI.Name, " by ", modI.Author, " from ", modI.RunURL)

		for k := range modI.Eval.Pre {
			evalString(modI.Eval.Pre[k])
		}

		evalUrl(modI.RunURL)

	} else {
		print("--> Loading ", modI.Name, " by ", modI.Author, " from ", censorPath(modI.baseDir))

		for k := range modI.Eval.Pre {
			evalString(modI.Eval.Pre[k])
		}
		if indexjsPath := filepath.Join(modI.baseDir, "index.js"); fileExists(indexjsPath) {
			requireFile(modI.baseName, indexjsPath)
		} else if mainjsPath := filepath.Join(modI.baseDir, "main.js"); fileExists(mainjsPath) {
			requireFile(modI.baseName, mainjsPath)
		} else {
			return errors.New("main.js and index.js do not exist")
		}

	}

	for k := range modI.Eval.Post {
		evalString(modI.Eval.Post[k])
	}

	print("--> Finished loading ", modI.Name, " by ", modI.Author, " from ", censorPath(modI.baseDir))

	return nil
}
