package common

import (
	"github.com/pkg/errors"
	"sync"
)

var Mode CommonMode
var modeLock = sync.Mutex{}

type CommonMode byte

const (
	CommonModeNo CommonMode = iota
	CommonModeJS
	CommonModeOS
)

var (
	ErrOSMode = errors.New("OS Mode Enabled")
	ErrJSMode = errors.New("JS Mode Enabled")
	ErrNoMode = errors.New("No Mode Enabled")
)

func LockAndCheckMode(mode CommonMode) error {
	modeLock.Lock()
	defer modeLock.Unlock()
	if Mode != CommonModeNo {
		return errors.New("Mode already set")
	}
	Mode = mode
	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
