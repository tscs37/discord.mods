#!/bin/sh
set -eu
set -o pipefail

# Verify we're in the installer director
cd ./assets/ && cd ../.. && cd ./installer

echo " * bootstrap.js"
echo "   - Installing bootstrap.js"
cp ./bootstrap.js ./assets/bootstrap.js

echo " * core.js "
echo "   - Building core.js"
cd ../d.mods/
./rebuild.sh
cd ../installer
echo "   - Installing Core.js"
cp ../d.mods/core.js ./assets/core.js

echo " * Discord.Mods API"
echo "   - Building Discord.Mods API"
cd ../d.mods.api
./rebuild.sh
cd ../installer
echo "   - Installing Discord.Mods API"
mkdir -p ./assets/mods/dmodsapi
cp ../d.mods.api/dmodapi.js ./assets/mods/dmodsapi/index.js
cp ../d.mods.api/dmodapi.dmod ./assets/mods/dmodsapi.dmod

echo " * 24h Timestamps"
echo "   - Installing 24h Timestamps"
mkdir -p ./assets/mods/24h-stamps
cp ../24h-stamps/index.js ./assets/mods/24h-stamps/index.js
cp ../24h-stamps/24h-stamps.dmod ./assets/mods/24h-stamps.dmod

echo Embedding Assets
rice embed-go -i .

echo Building Binary
go build
