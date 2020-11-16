const chatDiv = document.getElementById("chat-div");
chatDiv.style.height = window.innerHeight + "px";
const frame = document.getElementById("convo-frame");

function getBot() {
  console.log("bot");
}

function getConvo() {
  frame.setAttribute("src", "/");
}
