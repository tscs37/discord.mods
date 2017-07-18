package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var defaultDiscordPaths = []string{
	// Default path assumes discord is installed on the "discord" directory
	"./discord",
}

var resourcePath = []string{"resources", "app.asar"}

const (
	lookFor = "global.mainWindowId = mainWindow.id;"
	beginJs = `// BOOTSTRAP.JS HAS BEEN INSTALLED
// ABORT WHEN READING THIS MESSAGE`
	endJs = "// BOOTSTRAP.JS END HERE"
)

func getResourcePath(path string) string {
	return filepath.Join(append([]string{path}, resourcePath...)...)
}

var runner func(string, string) error = runIntAsar

var (
	flagUseExtAsar       = flag.Bool("ext-asar", true, "If disabled, will attempt to use the internal packer, not recommended")
	flagKillVanilla      = flag.Bool("overwrite", false, "If set, already installed discord.mod instances will be ignored")
	flagRestoreVanilla   = flag.Bool("restore", false, "If set, the vanilla app.asar will be copied back from app.asar.old")
	flagKillBackup       = flag.Bool("force-backup", false, "If set, a new app.asar.old is always created")
	flagAttemptReinstall = flag.Bool("reinstall", false, "If set, the installer attempts to clear previous bootstrap.js installs")
	flagOnlyCoreMods     = flag.Bool("only-mods", false, "Only copy mods and core.js files but do not reinstall the bootstrapper")
)

func dynamicPaths() []string {
	var ret = []string{}
	if runtime.GOOS == "darwin" {
		fmt.Println("Adding Darwin/MACOS Paths...")
		home := os.Getenv("HOME")
		if home != "" {
			ret = append(ret, filepath.Join(home, "Library", "Preference", "discord"))
		}
	} else if runtime.GOOS == "windows" {
		fmt.Println("Adding Windows Paths...")
		appdata := os.Getenv("APPDATA")
		if appdata != "" {
			ret = append(ret, filepath.Join(appdata, "discord"))
		}
	} else if strings.TrimSpace(runtime.GOOS) == "linux" {
		fmt.Println("Adding Linux Paths...")
		// Arch AUR Path
		ret = append(ret, filepath.Join("/opt", "discord"))
	} else {
		fmt.Println("Unsupported OS, you're probably going to have to enter your install path manually")
	}
	return ret
}

func main() {
	defaultDiscordPaths = append(defaultDiscordPaths, dynamicPaths()...)
	flag.Parse()
	if *flagUseExtAsar {
		runner = runExtAsar
	}
	defer fmt.Println()

	box := rice.MustFindBox("assets")
	base, err := getBase()
	if err != nil {
		fmt.Printf("Could not determine .discord.mods path: %s", err)
		return
	}

	fmt.Println("Checking Discord.Mods environment...")

	if _, err := os.Stat(base); os.IsNotExist(err) {
		fmt.Println("Creating folders...")
		if err = os.MkdirAll(base, 0755); err != nil {
			fmt.Println("COuld not make .discord.mods folder: ", err)
			return
		}
		if err = os.MkdirAll(filepath.Join(base, "mods"), 0755); err != nil {
			fmt.Println("Could not make .discord.mods/mods folder: ", err)
			return
		}
	}

	fmt.Println("Installing core.js...")
	err = ioutil.WriteFile(filepath.Join(base, "core.js"), box.MustBytes("core.js"), 0644)
	if err != nil {
		fmt.Println("Could not install core.js: ", err)
		return
	}

	if err := installMods(base, box); err != nil {
		fmt.Printf("Error while installing mods: %s", err)
		return
	}

	if *flagOnlyCoreMods {
		fmt.Println("Stopping install due to --only-mods")
		return
	}

	path, err := findDiscordPath()
	if err != nil {
		fmt.Printf("Error trying to read path: %s\n", err)
		return
	}

	bootstrapper := box.MustString("bootstrap.js")

	var hasBackupFile = false
	if _, err := os.Stat("app.asar.old"); !os.IsNotExist(err) {
		hasBackupFile = true
		if !(*flagKillVanilla || *flagRestoreVanilla) {
			fmt.Println("Found backup of app.asar.old, aborting until --restore or --overwrite it set")
			return
		} else if *flagKillVanilla && !*flagRestoreVanilla {
			fmt.Println("Ignoring existing app.asar.old")
		} else if *flagRestoreVanilla && !*flagKillVanilla {
			fmt.Println("Restoring vanilla discord...")
			err = copyFile(getResourcePath(path), "app.asar.old")
			if err != nil {
				fmt.Printf("Error on restore: %s\n", err)
				return
			}
			fmt.Println("Restore complete.")
		} else {
			fmt.Println("Either kill it or not, decide.")
			return
		}
	} else {
		fmt.Println("No backup file, making one later")
	}

	fmt.Printf("Installing Discord.Mods to %s\n", path)

	err = runner(path, bootstrapper)
	if err != nil {
		fmt.Printf("Error in ASAR engine: %s\n", err)
		return
	}

	fmt.Println("Encoding finished, copying files...")

	failBanner := `
	The installer has failed but you can complete the installation manually.

	Please backup your file at %s to another location
	and replace it with the output.asar in the current directory.

`

	if !hasBackupFile || *flagKillBackup {
		err = copyFile("app.asar.old", getResourcePath(path))
		if err != nil {
			fmt.Printf("Could not backup old ASAR: %s\n", err)
			fmt.Printf(failBanner, getResourcePath(path))
			return
		}
	} else {
		fmt.Println("There already is a backup file, not making a new one...")
	}
	err = copyFile(getResourcePath(path), "output.asar")
	if err != nil {
		fmt.Printf("Could not replace with new ASAR: %s\n", err)
		fmt.Printf(failBanner, getResourcePath(path))
		return
	}

	fmt.Println("Finished installation... Please restart Discord\n")
}

func findDiscordPath() (string, error) {
	for _, v := range defaultDiscordPaths {
		if _, err := os.Stat(filepath.Join(v, "resources", "app.asar")); os.IsNotExist(err) {
			continue
		}
		return v, nil
	}
	fmt.Println("Could not find Discord Installation Path")
	fmt.Print("Enter Path: >")
	for {
		bufin := bufio.NewReader(os.Stdin)
		line, err := bufin.ReadString('\n')

		line = strings.TrimRight(line, "\n")

		if err != nil {
			return "", errors.Wrap(err, "Error trying to read")
		}
		if _, err := os.Stat(filepath.Join(line, "resources", "app.asar")); os.IsNotExist(err) {
			fmt.Println("Path does not exist, trying again")
			continue
		}
		return line, nil
	}
}

func getBase() (string, error) {
	cuser, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Could not get current user")
	}
	return filepath.Join(cuser.HomeDir, ".discord.mods"), nil
}

func copyFile(dst, org string) error {
	stat, err := os.Stat(org)
	if err != nil {
		return errors.Wrap(err, "Could not stat Target File")
	}
	if stat.Size() == 0 {
		return errors.Errorf("Refusing to copy 0-byte file: %s", org)
	}
	file, err := os.OpenFile(org, os.O_RDONLY, stat.Mode())
	if err != nil {
		return errors.Wrap(err, "Could not open Target File")
	}
	defer file.Close()
	dest, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, stat.Mode())
	if err != nil {
		return errors.Wrap(err, "Could not open Destination File")
	}
	defer dest.Close()
	// Test if file is writable
	_, err = dest.WriteString("t")
	if err != nil {
		return errors.Wrap(err, "Write test failed")
	}
	err = dest.Truncate(0)
	if err != nil {
		// We fucked up, so be loud
		fmt.Printf(">> FILE %s WAS TRUNCATE AFTER WRITE <<", dst)
		fmt.Println(">> PLEASE CHECK IF THE FILE IS CORRUPTED <<")
		return errors.Wrap(err, "Truncate failed, File is corrupted")
	}
	_, err = dest.Seek(0, os.SEEK_SET)
	if err != nil {
		return errors.Wrap(err, "Could not seek start of file")
	}
	_, err = io.Copy(dest, file)
	if err != nil {
		return errors.Wrap(err, "Could not copy file manually")
	}
	return nil
}
