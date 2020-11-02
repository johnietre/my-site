const IP = "ws://129.119.172.61";
const PORT = ":8008/bot";
/*
 * Special query for bot that should be gotten using window.location.href
*/
var field = document.getElementById("field");
var ws = new WebSocket(IP+PORT, [], true);
ws.onopen = function() {
  ws.send("user");
};
ws.onerror = function(err) {
  console.log(err)
  document.getElementById("header").innerHTML = "SHIT";
};
ws.onmessage = function(msg) {
  console.log(msg.data);
};

function send() {
  msg = field.value;
  field.value = "";
  ws.send(msg);
};