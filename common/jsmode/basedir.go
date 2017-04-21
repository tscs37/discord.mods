package jsmode

import (
	"github.com/gopherjs/gopherjs/js"
)

var dmod = func() *js.Object { return js.Global.Get("dmodsNS") }
var global = func() *js.Object { return js.Global }

func GetHome() (string, error) {
	return dmod().Call("homedir").String(), nil
}
