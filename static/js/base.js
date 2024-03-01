document.addEventListener("keydown", (event) => {
  if (event.isComposing || document.activeElement != document.body) {
    return;
  }
  const searchDiv = document.querySelector("#base-search-div");
  if (event.key == "s") {
    searchDiv.hidden = false;
    searchDiv.querySelector("h3").innerText = "Search";
    searchDiv.querySelector("div[contenteditable]").focus();
    event.preventDefault();
  } else if (event.key == "S") {
    searchDiv.hidden = false;
    searchDiv.querySelector("h3").innerText = "Advanced Search";
    searchDiv.querySelector("div[contenteditable]").focus();
    event.preventDefault();
  } else if (event.key == "Escape") {
    searchDiv.hidden = true;
  }
});

document.querySelector("#base-search-div > div[contenteditable]")
    .addEventListener("keydown", (ev) => {ev.stopPropagation()});

document.querySelector("#base-search-cancel").addEventListener("click", () => {
  document.querySelector("#base-search-div").hidden = true;
});
