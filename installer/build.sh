#!/bin/sh

set -eu

. ./buildhelp.sh

mkBuildFolder

beginBuild

# Compile Core Assets
jsInstall       "bootstrap.js"              bootstrap
coreInstall     "core.js"                   core
modGoInstall    "Discord.Mods API"          dmodsapi

# Install Additional Mods
modJSInstall    "24h Timestamps"            24h-stamps
modGoInstall    "Discord.Mods CSS Loader"   dcss

endBuild

embed

finishBuild