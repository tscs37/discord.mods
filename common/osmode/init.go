package osmode

import (
	"go.rls.moe/misc/discord.mods/common"
	"path/filepath"
)

func init() {
	common.Must(Register())
}

func Register() error {
	if err := common.LockAndCheckMode(common.CommonModeOS); err != nil {
		// The common lib already has a mode so someone imported the wrong statement
		panic(err)
	}
	common.GetHome = GetHome
	common.GetFile = GetFile
	common.WriteFile = WriteFile
	common.Mkdir = Mkdir
	common.Exists = Exists
	common.JoinPath = filepath.Join
	return nil
}
