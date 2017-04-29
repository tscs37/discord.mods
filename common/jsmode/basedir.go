package jsmode

import (
	"github.com/gopherjs/gopherjs/js"
	"go.rls.moe/misc/discord.mods/common"
)

var dmod = func() *js.Object { return js.Global.Get("dmodsNS") }
var global = func() *js.Object { return js.Global }

// Returns user home or error
func GetHome() (string, error) {
	// We patch in a fix into the common code since the bootstrap.js
	// function seems to be broken but an update is not required
	dmodHomedirPatchlevel := dmod().Get("homedir_patch").Int()
	if dmodHomedirPatchlevel < patchLevels["homedir"] {
		println("Patching homedir() from patchlevel", dmodHomedirPatchlevel, "to", patchLevels["homedir"])
		dmod().Set("homedir", func() string {
			os, err := common.GetModule("os")
			if err != nil {
				println(err)
				return ""
			}
			return os.Call("homedir").String()
		})
		dmod().Set("homedir_patch", patchLevels["homedir"])
	}
	return dmod().Call("homedir").String(), nil
}
