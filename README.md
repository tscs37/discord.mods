# Discord.Mods

Discord.Mods is an unofficial modloader for Discord with a focus on reducing the
modding surface in the client and offloading as much as possible into external code.

## Installation

Atm there is no stable node.js-free installer, only an experimental one.

To install, simply compile the binary with Go of atleast version 1.8.

Running the binary is simple;

```bash
./installer
```

The installer will attempt to autodetect most configurations and install itself.

By default it uses the external ASAR engine, for this to work you need to have node.js with
a compatible `asar` binary in your path. If you do not, the installation will fail safely.

To use the internal engine use `--ext-asar=false`

### Installing ASAR

To install the asar binary, simply type;

```bash
npm -g install asar
```

## Updates

If the core.js or other files included in the installer have updated, simply run;

```bash
./installer --only-mods
```

### bootstrap.js Reinstall or Update

The installer uses multiple flags to ensure you won't accidentally kill your install.

The `--restore` flag instructs the installer to copy it's backup of the original Discord
back to the original location and install again

The `--overwrite` flag instructs the isntaller to ignore any pre-existing installs. This option **should only be used along with `--reinstall` and is not recommended unless you lost the original file**

If you changed the app.asar file in your Discord install, it's recommended to use `--force-backup`

The flag `--reinstall` instructs the installer to remove pre-existing bootstrap.js instances and reinsert them. This function is most likely less reliable than a simple `--reinstall` flag as it requires the `--overwrite` flag

## Notes on the installer

The installer is designed to operate as safely as possible.

It will not install to your installation if it detects a backup file in it's folder, assuming that this means that the bootstrap.js file is already installed.

If you already installed the bootstrap.js file, you can use `--only-mods` to update
the core libraries.

## Function

Discord.Mods uses the bootstrap.js file which looks for the file `~/.discord.mods/core.js`
or `%USER%/.discord.mods/core.js` on Windows, reads the contents and evals it.

`core.js` is a GopherJS application which then executes all further code. It also sets
up a namespace for public API functions; `dmodsNS`.

Mods are loaded from `.discord.mdos/mods`, each `.dmod` file represents a mod which must
be accompanied by a folder named like the file without extension.

Example; `24h-stamps.dmod` must put it's contents into `24h-stamps/`

The `.dmod` file is a simple yaml file that contains various information about the mod
including update urls and versioning.

## Setup Callback

Due to the simplicity of the `core.js` file, mods cannot have a dependency order.

Instead, when all mods are setup, the `core.js` file will execute a special
callback handler which notifies all other mods.

To register for this callback, call `dmodsNS.loadFinishedCallbackRegister`.

The only parameter is the function you wish to execute.

Ordering and Dependency Resolution is a task of a higher level mod.

## d.mods API

Discord.Mods comes with a simple notification API for changes in the browser
window plus an event handler.

To register for an event call `dmodsNS.onEvent(id, event, callback)`.

The ID must be unique to the mod, conflicting IDs will overwrite eachothers
callbacks.

To dispatch an event, call `dmodsNS.dispatchEvent(event, param)`

Param must be a single Javascript object.