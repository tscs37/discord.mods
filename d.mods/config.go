package main

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	baseFolder   = "/.discord.mods"
	configFile   = "/config.bin"
	modFolder    = "/mods"
	modExtension = ".dmod"
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
	base := getBase()
	if !dirExists(base) {
		mkdirPath(base, 0700)
	}
	cStr, err := getFile(filepath.Join(base, configFile))
	if err != nil {
		return nil, err
	}
	cRaw, err := base64.RawStdEncoding.DecodeString(cStr)
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
	base := getBase()
	if !dirExists(filepath.Join(base, modFolder)) {
		return nil, errors.New("Mod folder does not exist, check your setup")
	}
	cStr, err := getFile(filepath.Join(base, modFolder, modname))
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
	mod.baseDir = strings.TrimSuffix(filepath.Join(base, modFolder, modname), ".dmod")
	mod.baseName = modname
	return mod, nil
}

func saveConfig(c *config) error {
	cRaw, err := msgpack.Marshal(c)
	if err != nil {
		return err
	}
	home := os.Getenv("HOME")
	if home == "" {
		return errors.New("HOME path not set")
	}
	cStr := base64.RawStdEncoding.EncodeToString(cRaw)
	writeFile(filepath.Join(getBase(), configFile), cStr, 0600)
	return nil
}

func listMods() ([]string, error) {
	base := getBase()

	if !dirExists(filepath.Join(base, modFolder)) {
		mkdirPath(filepath.Join(base, modFolder), 0700)
		return []string{}, nil
	}
	files, err := lsDir(filepath.Join(base, modFolder))
	if err != nil {
		return nil, err
	}
	var mods = []string{}
	for _, v := range files {
		if strings.HasSuffix(v, modExtension) {
			mods = append(mods, v)
		}
	}
	return mods, nil
}
