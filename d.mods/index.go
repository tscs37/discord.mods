package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"sync"
	"time"
)

func main() {
	print("Testing JQuery Version")
	dmods.Call("loadJquery")
	// Wait for Jquery to become available
	time.Sleep(1 * time.Second)
	print("Version: ", getJqueryReliable())
	if getJqueryReliable() == "undefined" {
		print("No JQuery installed, check installation")
		return
	} else {
		print("Found JQuery: ", jquery.JQ)
		if getJqueryReliable() != "3.2.1" {
			print("Jquery version may be incompatible")
		}
	}
	print("Registering LoadFinished Callback")
	dmods.Set("loadFinishedCallbackRegister", loadFinishedCallbackRegister)
	print("Loading local keyring")
	config, err := loadConfig()
	if err != nil {
		print("Error loading keyring: ", err.Error())
		alert("Config is empty, using defaults!")
		print("Continuing with empty keyring...")
	}
	err = saveConfig(config)
	if err != nil {
		print("Error writing keyring")
		return
	}
	js.Global.Set("evalDModFile", evalFile)

	list, err := listMods()
	if err != nil {
		print("Error listing mods: ", err)
		return
	}
	print("Found ", len(list), " mods")
	wg := sync.WaitGroup{}
	for _, v := range list {
		mod, err := loadMod(v)
		if err != nil {
			print("Could not load mod '" + v + "': " + err.Error())
			continue
		}
		wg.Add(1)
		go func() {
			if err := evalMod(mod); err != nil {
				print("ERR: ", err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()
	print("Discord.Mods finished, notifying mods.")
	for k := range loadFinishedCbs {
		loadFinishedCbs[k].Invoke()
	}
	dmods.Delete("loadFinishedCallbackRegister")
	return
}

var loadFinishedCbs []*js.Object

func loadFinishedCallbackRegister(cb *js.Object) {
	if loadFinishedCbs == nil {
		loadFinishedCbs = []*js.Object{cb}
		return
	}
	loadFinishedCbs = append(loadFinishedCbs, cb)
}
