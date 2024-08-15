htmx.config.scrollIntoViewOnBoost = false;

//Browse bar ---------------------
function clear_browse_bar() {
  const browse_bar = document.getElementById("browseBar");
  browse_bar.value = "";
}
function cleanURL() {
  const splitURL = window.location.href.split("?q=");
  if (splitURL.length === 2 && splitURL[1] === "") {
    history.replaceState(null, null, splitURL[0]);
  }
}

function toggle_serie_switch() {
  var serie_switch = document.getElementById("serieModeToggle");
  var url = window.location.href;
  serie_switch.setAttribute("hx-get", url);
  // @ts-ignore
  htmx.process(serie_switch);
  // @ts-ignore
  htmx.trigger("#serieModeToggle", "toggle_serie_switch");
}

document.addEventListener("htmx:configRequest", function (configEvent) {
  const serie_switch = document.getElementById("serieModeToggle");
  // @ts-ignore
  configEvent.detail.headers["Smode"] = serie_switch.checked;
});

let access_token = null;
let refresh_token = null;
window.addEventListener("load", function (event) {
  const nat = new URLSearchParams(window.location.search).get("access-token");
  const nrt = new URLSearchParams(window.location.search).get("refresh-token");
  if (nat == null || nrt == null) {
    return;
  }
  access_token = nat;
  refresh_token = nrt;
  cleanURL();
});

function showTop() {
  window.scrollTo({ top: 0, behavior: "smooth" });
}
