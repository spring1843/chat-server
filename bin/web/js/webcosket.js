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

    $("#userInput").keypress(function(event) {
        if (event.which != 13) {
            return true;
        }
        event.preventDefault();
        if (!conn) {
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
        conn = new WebSocket(webSocketAddr);
        conn.onclose = function(evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = function(evt) {
            appendLog($("<div/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
});