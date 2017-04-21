package main

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"go.rls.moe/misc/discord.mods/common"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	keymap map[string]key    `msgpack:"keymap"`
	modmap map[string]string `msgpack:"modmap"`
}

type mod struct {
	Author        string      `yaml:"author"`     // Name of the Author
	Name          string      `yaml:"name"`       // Name of the Mod
	VersionNumber int         `yaml:"version"`    // Version Number does not represent an actual number, rather it's an absolute value of the current version, if not present updates must be done manually
	VersionURL    string      `yaml:"versionUrl"` // Download the current Version Number from here
	UpdateURL     string      `yaml:"updateUrl"`  // Template for downloading the new version
	RunURL        string      `yaml:"runUrl"`     // URL to include as script tag, if present this will load the script from that website and ignores local files
	baseDir       string      // Base Directory from which to execute main.js or index.js, whichever is present
	baseName      string      // Base Name of directory
	Eval          evalOptions `yaml:"eval"`
}

type evalOptions struct {
	Pre  []string `yaml:"pre"`  // Execute before Eval
	Post []string `yaml:"post"` // Execute after Eval
}

type key ed25519.PublicKey

func loadConfig() (*config, error) {
	base, err := common.GetBase()
	if err != nil {
		return nil, err
	}
	if exists, err := common.Exists(base); !exists && err == nil {
		err := common.Mkdir(base, 0700)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	cStr, err := common.GetFile(filepath.Join(base, common.ConfigFile))
	if err != nil {
		return nil, err
	}
	cRaw, err := base64.RawStdEncoding.DecodeString(string(cStr))
	if err != nil {
		return nil, err
	}
	c := &config{}
	err = msgpack.Unmarshal(cRaw, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func loadMod(modname string) (*mod, error) {
	base, err := common.GetBase()
	if err != nil {
		return nil, err
	}
	if exists, err := common.Exists(filepath.Join(base, common.ModFolder)); !exists && err == nil {
		return nil, errors.New("Mod folder does not exist, check your setup")
	} else if err != nil {
		return nil, errors.Wrap(err, "Could not check if path exists")
	}
	cStr, err := common.GetFile(filepath.Join(base, common.ModFolder, modname))
	if err != nil {
		return nil, err
	}
	var mod = &mod{
		VersionNumber: -1,
	}
	err = yaml.Unmarshal([]byte(cStr), mod)
	if err != nil {
		return nil, err
	}
	mod.baseDir = strings.TrimSuffix(filepath.Join(base, common.ModFolder, modname), ".dmod")
	mod.baseName = modname
	return mod, nil
}

func saveConfig(c *config) error {
	base, err := common.GetBase()
	if err != nil {
		return err
	}
	cRaw, err := msgpack.Marshal(c)
	if err != nil {
		return err
	}
	home := os.Getenv("HOME")
	if home == "" {
		return errors.New("HOME path not set")
	}
	cStr := base64.RawStdEncoding.EncodeToString(cRaw)
	return common.WriteFile(filepath.Join(base, common.ConfigFile), []byte(cStr), 0600)
}

func listMods() ([]string, error) {
	base, err := common.GetBase()
	if err != nil {
		return nil, err
	}

	if exists, err := common.Exists(filepath.Join(base, common.ModFolder)); err != nil {
		return nil, err
	} else if !exists {
		return []string{}, common.Mkdir(filepath.Join(base, common.ModFolder), 0700)
	}
	files, err := lsDir(filepath.Join(base, common.ModFolder))
	if err != nil {
		return nil, err
	}
	var mods = []string{}
	for _, v := range files {
		if strings.HasSuffix(v, common.ModExtension) {
			mods = append(mods, v)
		}
	}
	return mods, nil
}
