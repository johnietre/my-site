var convoDiv = document.getElementById("convo-div");
convoDiv.style.height = window.innerHeight + "px";
const slugIndex = window.location.href.indexOf("/convo");
const IP = window.location.href.substring(0, slugIndex).replace("http", "ws") + "socket/";
const PORT = ":8000";
/*
 * Special query for bot that should be gotten using window.location.href
 */
var field = document.getElementById("field");
var ws = new WebSocket(IP+PORT, [], true);
ws.onopen = function() {
  console.log("user");
};

ws.onerror = function(err) {
  console.log(err)
  document.getElementById("header").innerHTML = "Oh no";
};

ws.onmessage = function(msg) {
  console.log(msg.data);
};

function send() {
  msg = field.value;
  field.value = "";
  ws.send(msg);
};
