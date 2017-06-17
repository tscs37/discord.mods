package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/pkg/errors"
)

func runDiffAsar(path, diff string) error {
	fmt.Println(`

	You are using the binary-diff ASAR engine,
	this requires you to have an *exact* version of app.asar

	The engine will fail with an error if this is not the case.

	`)

	getResourcePath(diff)
	return nil
}

func makeDiffAsar(path, diff string) error {
	fmt.Println(`

	This will generate a binary diff from a app.asar.old
	and the installed app.asar.

	THIS MODE REQUIRES A INTALLED D.MOD bootstrap.js!!!

	Ensure the following:

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

	fingerprint
}
