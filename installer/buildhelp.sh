#!/bin/sh

set -eu

beginBuild() {
    echo "# 1/3 Building D.Mods"
}

endBuild() {
    echo "# 1/3 Finished Building D.Mods"
}

jsInstall() {
    echo " * $1"
    echo "   - Installing $1"
    cp ./"$2".js ./assets/"$2".js
}

coreInstall() {
    echo " * $1"
    echo "   - Building $1"
    cd ../"$2"/
    ./rebuild.sh
    cd ../installer
    echo "   - Installing $1"
    cp ../"$2"/"$2".js ./assets/"$2".js
}

modJSInstall() {
    echo " * $1"
    echo "   - Installing $1"
    mkdir -p ./assets/mods/"$2"
    cp ../"$2"/index.js ./assets/mods/"$2"/index.js
    cp ../"$2"/"$2".dmod ./assets/mods/"$2".dmod
}

modGoInstall() {
    echo " * $1"
    echo "   - Building $1"
    cd ../"$2"
    ./rebuild.sh
    cd ../installer
    echo "   - Installing $1"
    mkdir -p ./assets/mods/"$2"
    cp ../"$2"/index.js ./assets/mods/"$2"/index.js
    cp ../"$2"/"$2".dmod ./assets/mods/"$2".dmod
}

embed() {
    echo "# 2/3 Embedding Assets"
    rice embed-go -v -i .
    echo "# 2/3 Finished Embedding Assets"
}

finishBuild() {
    echo "# 3/3 Building Installer"
    export GOOS=${GOOS:=linux}
    export GOARCH=${GOARCH:=amd64}
    echo "   - Building build/installer_${GOOS}_${GOARCH}"
    go build -v -o build/installer_${GOOS}_${GOARCH}
    echo "# 3/3 Finished Building Installer"
}