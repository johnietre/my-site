const IP = "localhost";
const PORT = ":8080";

// Get document elems
var canvas = document.getElementById("chart");
var ctx = canvas.getContext("2d");

var symbolInput = document.getElementById("symbol-input");
var symbolSpan = document.getElementById("symbol-span");
var priceSpan = document.getElementById("price-span");

var orderForm = document.getElementById("order-form");
var orderButton = document.getElementById("order-button");
var cancelButton = document.getElementById("cancel-order-button");

// Connect and setup sockets
const symSock = new WebSocket("ws://"+IP+PORT, [], true);
const orderSock = new WebSocket("ws://"+IP+PORT+"/order", [], true);

// Setup canvas dimensions
ctx.canvas.width = window.innerWidth * 0.75;
ctx.canvas.height = window.innerWidth * 0.30;

// Search for a symbol
function searchSymbol() {
  var sym = symbolInput.value.toUpperCase();
  stock = {"sym": sym};
  symSock.send(JSON.stringify(stock));
  symbolSpan.innerHTML = sym;
  symbolInput.value = sym;
}

// Graph the symbol
function graph(prices) {

  var price = prices[1];
  var old = prices[0];

  if (price > old) priceSpan.style = "color: green;";
  else if (price < old) priceSpan.style = "color: red;";
  else priceSpan.style = "color: black;";
}

function sendOrder() {
  if (orderForm.hidden) {
    orderForm.hidden = false;
    orderButton.innerHTML = "Send Order";
    cancelButton.hidden = false;
    return;
  }
  order = {
    "account_id": 1234,
    "password": "rj385637",
    "sym": "AAPL",
    "qty": 120,
    "side": "sell",
  };
  orderSock.send(JSON.stringify(order));
}

function cancelOrder() {
  orderForm.hidden = true;
  cancelButton.hidden = true;
  orderButton.innerHTML = "Create Order";
}
