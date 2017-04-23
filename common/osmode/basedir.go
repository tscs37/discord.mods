package osmode

import (
	"github.com/pkg/errors"
	"os/user"
)

func GetHome() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Could not get user")
	}
	return user.HomeDir, nil
}
