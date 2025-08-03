(function() {
// 0 = closed, 1 = open, 2 = open (advanced)
let searchOpenStatus = 0;

document.addEventListener("keydown", (event) => {
  if (event.isComposing || document.activeElement != document.body) {
    return;
  }
  if (event.key == "s") {
    openSearch();
    event.preventDefault();
  } else if (event.key == "S") {
    openSearch(true);
    event.preventDefault();
  } else if (event.key == "Escape") {
    closeSearch();
  }
});

function openSearch(advanced) {
  alert("Search coming soon!");
  return;
  const searchDiv = document.querySelector("#base-search-div");
  if (advanced) {
    searchDiv.hidden = false;
    searchDiv.querySelector("h3").innerText = "Advanced Search";
    searchDiv.querySelector("div[contenteditable]").focus();
    searchOpenStatus = 2;
    return;
  }
  searchDiv.hidden = false;
  searchDiv.querySelector("h3").innerText = "Search";
  searchDiv.querySelector("div[contenteditable]").focus();
  searchOpenStatus = 1;
}

function closeSearch() {
  document.querySelector("#base-search-div").hidden = true;
  searchOpenStatus = 0;
}

function cycleSearch() {
  if (searchOpenStatus === 0) {
    openSearch();
  } else if (searchOpenStatus === 1) {
    openSearch(true);
  } else {
    closeSearch();
  }
}

document.querySelector("#base-search-div > div[contenteditable]")
    .addEventListener("keydown", (ev) => {ev.stopPropagation()});

document.querySelector("#base-search-cancel").addEventListener("click", () => {
  document.querySelector("#base-search-div").hidden = true;
});

function toggleMiniPopup(setShowing) {
  const nmp = document.querySelector("#nav-mini-popup");
  const burgerOpen = document.querySelector("#nav-bar-burger-open");
  const burgerClose = document.querySelector("#nav-bar-burger-close");
  if (setShowing === undefined) {
    setShowing = nmp.hidden;
  }
  if (setShowing) {
    burgerOpen.classList.add("hidden");
    burgerClose.classList.remove("hidden");
    nmp.hidden = false;
  } else {
    burgerOpen.classList.remove("hidden");
    burgerClose.classList.add("hidden");
    nmp.hidden = true;
  }
  /*
  if (setShowing === undefined) {
    nmp.style.display = (nmp.style.display === "none") ? "flex" : "none";
    return;
  }
  nmp.style.display = (setShowing) ? "flex" : "none";
  */
}

function doNavMiniSearch() {
  const searchInput = document.querySelector("#nav-mini-search-input");
  const searchCancel = document.querySelector("#nav-mini-search-cancel");
  if (searchInput.hidden) {
    searchInput.hidden = false;
    searchCancel.hidden = false;
    return;
  }
  searchInput.hidden = true;
  searchCancel.hidden = true;
  console.log("searched for:", searchInput.value);
}

function cancelNavMiniSearch() {
  const searchInput = document.querySelector("#nav-mini-search-input");
  const searchCancel = document.querySelector("#nav-mini-search-cancel");
  searchInput.hidden = true;
  searchCancel.hidden = true;
  searchInput.value = "";
}

window.cycleSearch = cycleSearch;
window.toggleMiniPopup = toggleMiniPopup;
window.doNavMiniSearch = doNavMiniSearch;
window.cancelNavMiniSearch = cancelNavMiniSearch;
})()
