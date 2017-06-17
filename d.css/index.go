package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"

	"bytes"
	"go.rls.moe/misc/discord.mods/common"
	_ "go.rls.moe/misc/discord.mods/common/jsmode"
	"html/template"
)

var dmods = js.Global.Get("dmodsNS")

func main() {
	println("Waiting for Discord.Mods API...")
	dmods.Call("loadFinishedCallbackRegister", func() {
		go runCSSLoaderInit()
	})
}

var jq func(args ...interface{}) jquery.JQuery

var cssTemplate = `<style type="text/css" class="customloaded">
{{ . }}
</style>`

func runCSSLoaderInit() {
	println("Discord.Mods API finished, loading Discord.Mods CSS loader")
	jq = jquery.NewJQuery
	csstmpl, err := template.New("css-template").Parse(cssTemplate)
	if err != nil {
		println("Could not load CSS Template, check your maintainer")
		return
	}

	base, err := common.GetBase()
	if err != nil {
		println("Could not determine base filepath")
		return
	}

	if dmods.Get("settings_dialog") == js.Undefined {
		println("No settings dialog provider, using theme.css")

		themeCss := common.JoinPath(base, "theme.css")
		cssData, err := common.GetFile(themeCss)
		if err != nil {
			println("Could not read theme.css, aborting")
			return
		}

		println(string(cssData))

		buffer := bytes.NewBuffer([]byte{})
		if err := csstmpl.Execute(buffer, template.CSS(cssData)); err != nil {
			println("Error executing CSS Template, check your maintainer")
			return
		}

		println(buffer.String())
		jq("head").Append(buffer.String())
	}
}
