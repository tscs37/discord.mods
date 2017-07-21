#!/bin/sh

set -eu

. ./buildhelp.sh

beginBuild

# Compile Core Assets
jsInstall       "bootstrap.js"              bootstrap
coreInstall     "core.js"                   core
modGoInstall    "Discord.Mods API"          dmodsapi

# Install Additional Mods
modGoInstall    "Discord.Mods CSS Loader"   dcss
# 24 hour Timestamps are Deprecated
# modJSInstall    "24h Timestamps"            24h-stamps

endBuild

embed
finishBuild