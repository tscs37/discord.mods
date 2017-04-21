package common

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
)

var GetJquery = func() (string, error) {
	return "undefined", errors.Wrap(ErrNoMode, "GetJquery")
}

var RequireFile = func(varname, file string) error {
	return errors.Wrap(ErrNoMode, "RequireFile")
}

var EvalFile = func(file string) error {
	return errors.Wrap(ErrNoMode, "EvalFile")
}

var EvalString = func(stmt string) error {
	return errors.Wrap(ErrNoMode, "EvalString")
}

var EvalURL = func(url string) error {
	return errors.Wrap(ErrNoMode, "EvalURL")
}

var GetModule = func(module string) (*js.Object, error) {
	return nil, errors.Wrap(ErrNoMode, "GetModule")
}

var Alert = func(string) {
	return
}

func GetFS() (*js.Object, error) {
	return GetModule("fs")
}
