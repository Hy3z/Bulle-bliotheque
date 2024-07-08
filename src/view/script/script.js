//Browse bar ---------------------
function clear_browse_bar() {
    var browse_bar = document.getElementById("browse-bar");
    browse_bar.value = '';
}
function cleanURL() {
    var splitURL = window.location.href.split("?q=");
    if (splitURL.length == 2 && splitURL[1] == '') {
        history.replaceState(null, null, splitURL[0]);
    }
}
function toggle_serie_switch() {
    var serie_switch = document.getElementById("toggle_serie_mode");
    var url = window.location.href;
    serie_switch.setAttribute("hx-get", url);
    // @ts-ignore
    htmx.process(serie_switch);
    // @ts-ignore
    htmx.trigger("#toggle_serie_mode", "toggle_serie_switch");
}
document.addEventListener("htmx:configRequest", function (configEvent) {
    var serie_switch = document.getElementById("toggle_serie_mode");
    // @ts-ignore
    configEvent.detail.headers['Smode'] = serie_switch.checked;
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
