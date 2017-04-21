package osmode

import (
	"github.com/pkg/errors"
	"go.rls.moe/misc/discord.mods/common"
	"os/user"
	"io/ioutil"
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
	return nil
}

func GetHome() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Could not get user")
	}
	return user.HomeDir, nil
}


func GetFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}