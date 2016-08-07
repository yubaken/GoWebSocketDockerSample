var messageTxt;
var messages;

window.onload = function () {
    $.ajax({
        url: "/messages",
        cache: false,
        success: function (data) {
            var messages = data.messages;
            for (var i = 0; i < messages.length; i++) {
                appendMessage($("<div>" + messages[i] + "</div>"));
            }
        }
    });
};

$(function () {
    messageTxt = $("#messageTxt");
    messages = $("#messages");
    var msg = {Room: "chat1", Msg: ""};

    ws = new Ws("ws://" + HOST + "/my_endpoint");
    ws.OnConnect(function () {
        console.log("Websocket connection established");
        ws.Emit("init", JSON.stringify(msg));
    });

    ws.OnDisconnect(function () {
        ws.Emit("leave", JSON.stringify(msg));
        appendMessage($("<div class='center'><h3>Disconnected</h3></div>"));
    });

    ws.On("join", function (message) {
        prependMessage($("<div>" + message + "</div>"));
    });

    ws.On("chat", function (message) {
        appendMessage($("<div>" + message + "</div>"));
    });

    $("#sendBtn").click(function () {
        msg.Msg = messageTxt.val().toString();
        ws.Emit("chat", JSON.stringify(msg));
        messageTxt.val("");
    });
});

function prependMessage(messageDiv) {
    messages.prepend(messageDiv);
}

function appendMessage(messageDiv) {
    var theDiv = messages[0];
    var doScroll = theDiv.scrollTop == theDiv.scrollHeight - theDiv.clientHeight;
    messageDiv.appendTo(messages);
    if (doScroll) {
        theDiv.scrollTop = theDiv.scrollHeight - theDiv.clientHeight;
    }
}