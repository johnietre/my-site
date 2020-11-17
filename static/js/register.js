var password = document.getElementById("password-input").value;
var confirmation = document.getElementById("confirmation-input").value;
var errp = document.getElementById("error-p");

if (errp.innerHTML != "") errp.hidden = false;

document.getElementById("submit-input").addEventListener("click", function(event) {
  errp.hidden = true;
  if (password != confirmation || password.length < 8) {
    if (password.length < 8) errp.innerHTML = "Password must be at least 8 characters";
    else errp.innerHTML = "Password and confirmation must match!";
    errp.hidden = false;
    event.preventDefault();
  }
});