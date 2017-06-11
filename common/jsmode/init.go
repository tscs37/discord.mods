package jsmode

import "go.rls.moe/misc/discord.mods/common"

func init() {
	common.Must(Register())
}

func Register() error {
	if err := common.LockAndCheckMode(common.CommonModeJS); err != nil {
		// The common lib already has a mode so someone imported the wrong statement
		panic(err)
	}
	common.GetHome = GetHome
	common.GetJquery = GetJquery
	common.RequireFile = RequireFile
	common.Exists = Exists
	common.Mkdir = Mkdir
	common.GetFile = GetFile
	common.WriteFile = WriteFile
	common.Alert = Alert
	common.EvalFile = EvalFile
	common.EvalString = EvalString
	common.EvalURL = EvalURL
	common.GetModule = GetModule
	common.JoinPath = JoinPath
	return nil
}

var patchLevels = map[string]int{
	"homedir": 1,
}
