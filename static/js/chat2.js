var messageTxt;
var messages;

$(function () {
    messageTxt = $("#messageTxt");
    messages = $("#messages");
    var msg = {Room: "chat2", Msg: ""};

    ws = new Ws("ws://" + HOST + "/my_endpoint");
    ws.OnConnect(function () {
        console.log("Websocket connection established");
        ws.Emit("init", JSON.stringify(msg));
    });

    ws.OnDisconnect(function () {
        ws.Emit("leave", JSON.stringify(msg));
        appendMessage($("<div><center><h3>Disconnected</h3></center></div>"));
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

function appendMessage(messageDiv) {
    var theDiv = messages[0];
    var doScroll = theDiv.scrollTop == theDiv.scrollHeight - theDiv.clientHeight;
    messageDiv.appendTo(messages);
    if (doScroll) {
        theDiv.scrollTop = theDiv.scrollHeight - theDiv.clientHeight;
    }
}