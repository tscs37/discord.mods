function install_24hstamps() {
    console.log("Finishing 24h Timestamp Install");
    var convert = function() {
        $(".timestamp").each(function() {
            var t = $(this);

            if(t.data("24") != undefined) return;

            var text = t.text();
            var matches = /(.*)?at\s+(\d{1,2}):(\d{1,2})\s+(.*)/.exec(text);
            if(matches == null) return false;
            if(matches.length < 5) return false;

            var h = parseInt(matches[2]);
            if(matches[4] == "AM") {
                if(h == 12) h -= 12;
            }else if(matches[4] == "PM") {
                if(h < 12) h += 12;
            }

            matches[2] = ('0' + h).slice(-2);
            t.text(matches[1] + matches[2] + ":" + matches[3]);
            t.data("24", true);
        });
    };
    dmodsNS.onEvent("24hourstamps", "newMessage", convert);
    dmodsNS.onEvent("24hourstamps", "channelSwitch", convert);
}

dmodsNS.loadFinishedCallbackRegister(install_24hstamps);