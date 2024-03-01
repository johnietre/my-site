//function appsSubmitReviewForm(form) {
document.querySelector("#apps-review-form").addEventListener("submit", (ev) => {
  const form = ev.target;
  if (form["name"].value === "" || form["reason"].value === "") {
    ev.preventDefault();
    alert("Must select app name");
    return;
  } else if (form["description"].value.length > maxDescriptionLen) {
    ev.preventDefault();
    alert("Description too long");
    return;
  }
  if (!confirm("Are you sure you want to submit?")) {
    ev.preventDefault();
  }
  //appsClearForm(form);
});

document.querySelector("#apps-review-form").addEventListener("htmx:beforeSwap", (ev) => {
    if (ev.detail.shouldSwap) {
      //document.querySelector("#apps-form-result").style.color = "";
      ev.detail.target = htmx.find("#apps-form-result");
      ev.detail.shouldSwap = true;
      ev.detail.isError = false;
      ev.detail.target.innerHTML = ev.detail.xhr.response;
      appsClearForm(ev.detail.elt);
      alert("Success");
    } else {
      // TODO: Handle better
      ev.detail.target = htmx.find("#apps-form-result");
      ev.detail.shouldSwap = true;
      ev.detail.isError = false;
      ev.detail.target.innerHTML = ev.detail.xhr.response;
      alert(`Error`);
    }
});

function appsClearForm(form) {
  if (form === undefined) {
    form = document.querySelector("#apps-review-form");
  }
  for (const el of form.elements) {
    el.value = "";
  }
  appsUpdateDescCountSpan(0);
}

document.querySelector("#apps-review-form")["description"].addEventListener("input", (ev) => {
  appsUpdateDescCountSpan(ev.target.value.length);
});

function appsUpdateDescCountSpan(count) {
  if (count === undefined) {
    count = document.querySelector("#apps.review-form")["description"].value.length;
  }
  const countSpan = document.querySelector("#apps-desc-count");
  countSpan.innerText = count;
  countSpan.style.setProperty("color", (count > maxDescriptionLen) ? "red" : "");
}

function appsFocusFormWithApp(name) {
  const form = document.querySelector("#apps-review-form");
  form["app"].value = name;
  form["email"].focus();
}
