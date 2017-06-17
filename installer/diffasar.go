package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/chzyer/readline"
	"github.com/monmohan/xferspdy"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
	"io/ioutil"
	"os"
)

type diffStruct struct {
	Checksum []byte
	Diff     []xferspdy.Block
}

func runDiffAsar(path, diff string) error {
	fmt.Println(`

	You are using the binary-diff ASAR engine,
	this requires you to have an *exact* version of app.asar

	The engine will fail with an error if this is not the case.

	`)
	box := rice.MustFindBox("assets")
	var diffData = diffStruct{}
	diffRaw := box.MustBytes("app.asar.diff")
	err := msgpack.Unmarshal(diffRaw, &diffData)
	if err != nil {
		return errors.Wrap(err, "Could not unpack diff file")
	}
	// Todo
	return nil
}

func makeDiffAsar(path, diff string) error {
	fmt.Println(`

	This will generate a binary diff from a app.asar.old
	and the installed app.asar.

	THIS MODE REQUIRES A INTALLED D.MOD bootstrap.js!!!

	Ensure the following is present in the **CURRENT** folder:

		- app.asar with bootstrap.js patch
		- app.asar.old from the original app.asar

	`)

	prompt, err := readline.Line("Continue? (Y)es/(n)o: ")
	if err != nil {
		return errors.Wrap(err, "Could not read from Console")
	}
	if prompt != "Y" {
		return errors.New("Did not press Y, aborting")
	}

	fingerprint := xferspdy.NewFingerprint("./app.asar.old", 4096)

	binaryDiff := xferspdy.NewDiff("./app.asar", *fingerprint)

	checksum, err := hashFile("./app.asar")
	if err != nil {
		return errors.Wrap(err, "Could not hash app.asar")
	}

	diffOut := diffStruct{
		Diff:     binaryDiff,
		Checksum: checksum,
	}
	diffData, err := msgpack.Marshal(diffOut)
	if err != nil {
		return errors.Wrap(err, "Could not marshall diff data")
	}
	err = ioutil.WriteFile("app.asar.diff", diffData, 0666)
	if err != nil {
		return errors.Wrap(err, "Could not write diff")
	}
	return nil
}

func hashFile(file string) ([]byte, error) {
	hasher, err := blake2b.New512([]byte{})
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize hasher")
	}
	toHash, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open file")
	}
	defer toHash.Close()

	_, err = io.Copy(hasher, toHash)
	if err != nil {
		return nil, errors.Wrap(err, "Could not hash file")
	}

	return hasher.Sum(nil), nil
}
