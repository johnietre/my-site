// function productsSubmitReviewForm(form) {
document.querySelector("#products-review-form")
    .addEventListener("submit", (ev) => {
      const form = ev.target;
      if (form["name"].value === "" || form["reason"].value === "") {
        ev.preventDefault();
        alert("Must select product name");
        return;
      } else if (form["description"].value.length > maxDescriptionLen) {
        ev.preventDefault();
        alert("Description too long");
        return;
      }
      if (!confirm("Are you sure you want to submit?")) {
        ev.preventDefault();
      }
      // productsClearForm(form);
    });

document.querySelector("#products-review-form")
    .addEventListener("htmx:beforeSwap", (ev) => {
      if (ev.detail.shouldSwap) {
        // document.querySelector("#products-form-result").style.color = "";
        /*
        ev.detail.target = htmx.find("#products-form-result");
        ev.detail.shouldSwap = true;
        ev.detail.target.innerHTML = ev.detail.xhr.response;
        */
        ev.detail.shouldSwap = false;
        ev.detail.isError = false;
        productsClearForm(ev.detail.elt);
        alert("Success");
      } else {
        // TODO: Handle better
        ev.detail.target = htmx.find("#products-form-result");
        ev.detail.shouldSwap = true;
        ev.detail.isError = false;
        ev.detail.target.innerHTML = ev.detail.xhr.response;
        alert(`Error`);
      }
    });

function productsClearForm(form) {
  if (form === undefined) {
    form = document.querySelector("#products-review-form");
  }
  for (const el of form.elements) {
    el.value = "";
  }
  document.querySelector("#products-form-result").innerHTML = "";
  productsUpdateDescCountSpan(0);
}

document.querySelector("#products-review-form")["description"].addEventListener(
    "input", (ev) => { productsUpdateDescCountSpan(ev.target.value.length); });

function productsUpdateDescCountSpan(count) {
  if (count === undefined) {
    count = document.querySelector("#products.review-form")["description"]
                .value.length;
  }
  const countSpan = document.querySelector("#products-desc-count");
  countSpan.innerText = count;
  countSpan.style.setProperty("color",
                              (count > maxDescriptionLen) ? "red" : "");
}

function productsFocusFormWithProduct(name) {
  const form = document.querySelector("#products-review-form");
  form["product"].value = name;
  // form["email"].focus();
  form["platform"].focus();
}
