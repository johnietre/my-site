const IP = "129.119.172.61";
const PORT = ":8080";

// Get document elems
var canvas = document.getElementById("chart");
var ctx = canvas.getContext("2d");
var symbolInput = document.getElementById("symbol-input");
var symbolSpan = document.getElementById("symbol-span");
var priceSpan = document.getElementById("price-span");

// Connect and setup sockets
var symSock = new WebSocket(IP+PORT, [], true);
var orderSock = new WebSocket(IP+PORT+"/order", [], true);

// Setup canvas dimensions
ctx.canvas.width = window.innerWidth * 0.75;
ctx.canvas.height = window.innerWidth * 0.30;

// Search for a symbol
function search() {
  var sym = symbolInput.value.toUpperCase();
  stock = {"sym": sym};
  symSock.send(JSON.stringify(stock));
  symbolSpan.innerHTML = sym;
}

// Graph the symbol
function graph(prices) {

  var price = prices[1];
  var old = prices[0];

  if (price > old) priceSpan.style = "color: green;";
  else if (price < old) priceSpan.style = "color: red;";
  else priceSpan.style = "color: black;";
}
