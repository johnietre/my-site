(function () {
  // 0 = closed, 1 = open, 2 = open (advanced)
  let searchOpenStatus = 0;

  function setupListeners() {
    try {
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

      document
        .querySelector("#base-search-entry")
        .addEventListener("keydown", (ev) => {
          ev.stopPropagation();
        });

      document
        .querySelector("#base-search-cancel")
        .addEventListener("click", () => {
          document.querySelector("#base-search-div").hidden = true;
        });
    } catch (e) {
      console.error(e);
    }
  }

  function openSearch(advanced) {
    alert("Search coming soon!");
    return;
    const searchDiv = document.querySelector("#base-search-div");
    if (advanced) {
      searchDiv.hidden = false;
      searchDiv.querySelector("base-search-div-header").innerText =
        "Advanced Search";
      searchDiv.querySelector("base-search-entry").focus();
      searchOpenStatus = 2;
      return;
    }
    searchDiv.hidden = false;
    searchDiv.querySelector("base-search-div-header").innerText = "Search";
    searchDiv.querySelector("base-search-entry").focus();
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

  const LIGHT_ICON_NAME = "fa-sun-o";
  const DARK_ICON_NAME = "fa-moon-o";
  const SYSTEM_ICON_NAME = "fa-cog";

  function cycleTheme() {
    const toggleIcon = document.querySelector("#nav-bar-theme-toggle-icon");
    let theme = "dark";
    if (toggleIcon.classList.contains(LIGHT_ICON_NAME)) {
      theme = "system";
    } else if (toggleIcon.classList.contains(DARK_ICON_NAME)) {
      theme = "light";
    }
    setTheme(theme);
    saveTheme(theme);
  }

  function setTheme(theme) {
    const toggleIcon = document.querySelector("#nav-bar-theme-toggle-icon");
    const body = document.body;
    if (theme == "light") {
      toggleIcon.classList.remove(DARK_ICON_NAME);
      toggleIcon.classList.remove(SYSTEM_ICON_NAME);
      toggleIcon.classList.add(LIGHT_ICON_NAME);
      body.classList.remove("dark-theme");
      body.classList.add("light-theme");
    } else if (theme == "system") {
      toggleIcon.classList.remove(LIGHT_ICON_NAME);
      toggleIcon.classList.remove(DARK_ICON_NAME);
      toggleIcon.classList.add(SYSTEM_ICON_NAME);
      body.classList.remove("light-theme");
      body.classList.remove("dark-theme");
    } else {
      toggleIcon.classList.remove(SYSTEM_ICON_NAME);
      toggleIcon.classList.remove(LIGHT_ICON_NAME);
      toggleIcon.classList.add(DARK_ICON_NAME);
      body.classList.remove("light-theme");
      body.classList.add("dark-theme");
    }
  }

  const THEME_KEY = "johnietre-theme";

  function saveTheme(theme) {
    window.localStorage.setItem(THEME_KEY, theme);
  }

  function loadTheme() {
    const theme = window.localStorage.getItem(THEME_KEY);
    switch (theme) {
      case "light":
      case "system":
        return theme;
      default:
        return "dark";
    }
  }

  setupListeners();
  setTheme(loadTheme());

  window.cycleSearch = cycleSearch;
  window.toggleMiniPopup = toggleMiniPopup;
  window.doNavMiniSearch = doNavMiniSearch;
  window.cancelNavMiniSearch = cancelNavMiniSearch;
  window.cycleTheme = cycleTheme;
})();
