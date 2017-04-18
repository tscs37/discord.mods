// BOOTSTRAP.JS HAS BEEN INSTALLED
// ABORT WHEN READING THIS MESSAGE

function init_DMODS() {

    try {
        mainWindow.webContents.on('dom-ready', function () {
            mainWindow.webContents.executeJavaScript(
                `
function startDiscordMods()
{
    var fs = require("fs");
    var os = require("os");
    var base = os.homedir() + "/.discord.mods/core.js";

    fs.readFile(base, function read(err, data) {
        if (err) {
            throw err;
        }
        var content = data;
        if (Buffer.isBuffer(data)) {
            content = data.toString();
        }
        eval(content);
    });
}

function makeDmodsNamespace()
{
    return {
        getFS: function () {
            return require("fs");
        },
        loadJquery: function () {
            var jquery = document.createElement('script');
            jquery.setAttribute('src', 'https://code.jquery.com/jquery-3.2.1.min.js');
            document.head.appendChild(jquery);
            return jquery
        },
        loadRemote: function(url) {
            var scr = document.createElement('script');
            scr.setAttribute('src', url);
            document.head.appendChild(scr);
        },
        homedir: function() {
            return require('os').homedir();
        },


        nil: function(){}
    }
}

dmodsNS = makeDmodsNamespace();
dmodsNS.loadJquery();
startDiscordMods();
`);
        });

    } catch(err) {
        console.log(err);
    }

}

init_DMODS();

// BOOTSTRAP.JS END HERE