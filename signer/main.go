package main

import (
	"flag"
	"fmt"
	"go.rls.moe/misc/discord.mods/common"
	_ "go.rls.moe/misc/discord.mods/common/osmode"
)

var (
	flagFileTosign = flag.String("file", "mymod.dmod", "File to sign with your key")
	flagKeyFile    = flag.String("key-file", "mykey.key", "Key with which to sign, if not present it is generated")
	flagNoPasswd   = flag.Bool("no-passwd", false, "Disable Password Entry for keys without password")
)

func main() {
	flag.Parse()
	home, err := common.GetBase()
	if err != nil {
		return
	}
	fmt.Println(home)
}
