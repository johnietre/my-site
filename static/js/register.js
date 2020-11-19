var passwordInput = document.getElementById("password-input");
var confirmationInput = document.getElementById("confirmation-input");
var errp = document.getElementById("error-p");

if (errp.innerHTML != "") errp.hidden = false;

document.getElementById("submit-input").addEventListener("click", function(event) {
  var password = passwordInput.value;
  var confirmation = confirmationInput.value;
  errp.hidden = true;
  if (password != confirmation || password.length < 8) {
    if (password.length < 8) errp.innerHTML = "Password must be at least 8 characters";
    else errp.innerHTML = "Password and confirmation must match!";
    errp.hidden = false;
    event.preventDefault();
  }
});