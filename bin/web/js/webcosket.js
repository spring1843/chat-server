'use strict';

$(function() {
    var conn;
    var msg = $("#userInput");
    var log = $("#log");
    var hostAndPort = location.hostname+(location.port ? ':'+location.port: '');
    var webSocketAddr = "wss://" + hostAndPort + "/ws";

    var validCommands = ["run-function"];
    var timeStampLengh = 11

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
        if (msg.length < commandLength + timeStampLengh ){
            continue;
        }
        var part = msg.substr(timeStampLengh,commandLength);
        if (part == validCommands[i]) {
             isCommand = true;
             command = validCommands[i];
        }
    }

    if (!isCommand){
        return;
    }

    if (command == "run-function"){
        var parts = msg.match(/[^{}]+(?=\})/g);
        if (parts[0] == "setChannel"){
            setRoom(parts[1]);
        }
    }
    }

    function connect(host){
       conn = new WebSocket(host);
            conn.onclose = function(evt) {
                appendLog($("<div><b>Connection closed. Type anything and enter to reconnect.</b></div>"))
            }
            conn.onmessage = function(evt) {
                var dataString = evt.data;
                runIfCommand(dataString);
                appendLog($("<div/>").text(dataString.substr(timeStampLengh, dataString.length - timeStampLengh)))
            }
    }

    $("#userInput").keypress(function(event) {
        if (event.which != 13) {
            return true;
        }
        event.preventDefault();
        if (!conn || conn.readyState != 1) {
            connect(webSocketAddr);
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