//Browse bar ---------------------
function clear_browse_bar() {
  var browse_bar = document.getElementById("browseBar");
  browse_bar.value = "";
}
function cleanURL() {
  var splitURL = window.location.href.split("?q=");
  if (splitURL.length == 2 && splitURL[1] == "") {
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
  var serie_switch = document.getElementById("serieModeToggle");
  // @ts-ignore
  configEvent.detail.headers["Smode"] = serie_switch.checked;
});

var access_token = null;
var refresh_token = null;
window.addEventListener("load", function (event) {
  var nat = new URLSearchParams(window.location.search).get("access-token");
  var nrt = new URLSearchParams(window.location.search).get("refresh-token");
  if (nat == null || nrt == null) {
    return;
  }
  access_token = nat;
  refresh_token = nrt;
  cleanURL();
});

/*
var header = document.getElementById("scrollHeader");
var sticky = header.offsetTop;
window.addEventListener("scroll", function () {
    if (window.pageYOffset > sticky) {
        header.classList.add("sticky");
    } else {
        header.classList.remove("sticky");
    }
})
*/

function showTop() {
  window.scrollTo({ top: 0, behavior: "smooth" });
}
