package common

import (
	"path/filepath"
	"strings"
)

const (
	BaseFolder   = "/.discord.mods"
	ConfigFile   = "/config.bin"
	ModFolder    = "/mods"
	ModExtension = ".dmod"
	ModSignature = ".dmod.sig"
	SignerKey    = "signer.key"
)

var GetHome func() (string, error)

func GetBase() (string, error) {
	home, err := GetHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, BaseFolder), nil
}

func CensorPath(path string) (string, error) {
	home, err := GetHome()
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(path, home), nil
}
