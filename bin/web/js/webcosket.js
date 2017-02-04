'use strict';

$(function() {
    var conn;
    var msg = $("#userInput");
    var log = $("#log");
    var hostAndPort = location.hostname+(location.port ? ':'+location.port: '');
    var webSocketAddr = "wss://" + hostAndPort + "/ws";

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    function connect(host){
       conn = new WebSocket(host);
            conn.onclose = function(evt) {
                appendLog($("<div><b>Connection closed. Type anything and enter to reconnect.</b></div>"))
            }
            conn.onmessage = function(evt) {
                appendLog($("<div/>").text(evt.data))
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