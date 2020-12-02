const chatDiv = document.getElementById("chat-div");
const convoDiv = document.getElementById("convo-div");
const convosList = document.getElementById("convos-list");
const field = document.getElementById("field");
const userSearch = document.getElementById("user-search");
const convoHeader = document.getElementById("convo-header");

const slugIndex = window.location.href.indexOf("/convo");
const url = window.location.href.substring(0, slugIndex).replace("http", "ws") + "socket/";
const ws = new WebSocket(url, [], true);
var xhr = new XMLHttpRequest();

var current = "";

chatDiv.style.height = window.innerHeight + "px";
convoDiv.style.height = window.innerHeight + "px";

ws.onopen = function() {
  console.log("Connected");
};
ws.onerror = function() {
  console.log(err);
  document.getElementById("convo-header").innerHTML = "Oh no!";
};
ws.onmessage = function(msgJSON) {
  msg = JSON.parse(msgJSON.data)
  console.log(msg);
};
function send() {
  msg = field.value;
  field.value = "";
  // ws.send(msg);
}

xhr.onload = function(e) {
  if (CharacterData.readyState === 4) {
    if (xhr.status === 200) {
      messages = JSON.parse(xhr.responseText);
      console.log(messages)
    } else {
      console.error(xhr.statusText, e);
    }
  }
};
xhr.onerror = function(e) {
  console.error(xhr.statusText, e);
};

function newConvo() {
  if (userSearch.hidden) {
    userSearch.hidden = false;
    return;
  }
  var user = userSearch.value;
  userSearch.value = "";

  // Add auth (uname and password to open method)
  // xhr.open("GET", "/chat/convo", true);
  // xhr.send(null);

  for (var u of convosList.children) {
    var b = u.children[0];
    if (b.innerHTML == user) {
      userSearch.hidden = true;
      return;
    }
  }
  var button = document.createElement("button");
  button.setAttribute("onclick", "getConvo(this)");
  button.innerHTML = user;
  var li = document.createElement("li");
  li.appendChild(button);
  convosList.appendChild(li);
  convoHeader.innerHTML = user;
  userSearch.hidden = true;
}

function getConvo(caller) {
  // frame.setAttribute("src", "/");
  current = caller.innerHTML;
  convoHeader.innerHTML = current;

  // Add auth (uname and password to open method)
  // xhr.open("GET", "/chat/convo", true);
  // xhr.send(null);
}

function getBot() {
  convoHeader.innerHTML = "Bot";
}

/*
const frame = document.getElementById("convo-frame");
function getConvo() {
  frame.setAttribute("src", "/");
}
*/
