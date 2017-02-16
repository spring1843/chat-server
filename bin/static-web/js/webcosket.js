'use strict';

$(function() {
    var conn;
    var msg = $("#userInput");

    $(msg).focus();

    var log = $("#log");
    var hostAndPort = location.hostname+(location.port ? ':'+location.port: '');

    var protocol = "wss:"
    if (location.protocol != 'https:'){
      protocol = "ws:"
    }

    var webSocketAddr = protocol + "//" + hostAndPort + "/ws";

    var validCommands = ["06"]; // UserOutPutTypeFERunFunction
    var timeStampLentgh = 13

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    function runIfCommand(msg){
    var isCommand = false;
    var command = "";

    for (var i = 0; i < validCommands.length ; i++) {
        var commandLength = validCommands[i].length;
        if (msg.length < commandLength + timeStampLentgh ){
            continue;
        }
        var part = msg.substr(timeStampLentgh,commandLength);
        if (part == validCommands[i]) {
             isCommand = true;
             command = validCommands[i];
        }
    }

    if (!isCommand){
        return;
    }

    if (command == "06"){
        var parts = msg.match(/[^{}]+(?=\})/g);
        if (parts[0] == "setChannel"){
            setRoom(parts[1]);
        }
    }
    }

    function connect(host){
       conn = new WebSocket(host);
            conn.onclose = function(evt) {
                appendLog($("<div><b>Connection Error. Press Enter or return to try again.</b></div>"))
            }
            conn.onmessage = function(evt) {
                var dataString = evt.data;
                runIfCommand(dataString);
                appendLog($("<div/>").text(dataString.substr(timeStampLentgh, dataString.length - timeStampLentgh)))
            }
    }

    $("#userInput").keypress(function(event) {
        if (event.which != 13) {
            return true;
        }
        event.preventDefault();
        if (!conn || conn.readyState != 1) {
            connect(webSocketAddr);
            msg.val("");
            return false;
        }

        if (!msg.val()) {
            return false;
        }
        conn.send(msg.val());
        msg.val("");
        return false
    });

    if (window["WebSocket"]) {
        connect(webSocketAddr);
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
});
