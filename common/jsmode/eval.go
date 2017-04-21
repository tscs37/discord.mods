package jsmode

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
	"go.rls.moe/misc/discord.mods/common"
	"io/ioutil"
	"net/http"
)

func GetJquery() (string, error) {
	if global().Get("jQuery") == js.Undefined {
		return "undefined", errors.New("JQuery is undefined")
	}
	if global().Get("jQuery").Get("fn") == js.Undefined {
		return "undefined", errors.New("JQuery FN is undefined")
	}
	return global().Get("jQuery").Get("fn").Get("jquery").String(), nil
}

func RequireFile(varname, file string) error {
	dmod().Set(varname, global().Call("require", file))
	return nil
}

func EvalFile(file string) error {
	dat, err := common.GetFile(file)
	if err != nil {
		return err
	}
	return common.EvalString(string(dat))
}

func EvalString(stmt string) error {
	global().Call("eval", stmt)
	return nil
}

func EvalURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Error loading resource")
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Error reading resource")
	}
	return common.EvalString(string(dat))
}

func GetModule(module string) (*js.Object, error) {
	return global().Call("require", module), nil
}

func Alert(msg string) {
	global().Call("alert", msg)
}
