//Browse bar ---------------------
function clear_browse_bar() {
    let browse_bar = <HTMLInputElement>document.getElementById("browse-bar");
    browse_bar.value = '';
}

function cleanURL() {
    let splitURL = window.location.href.split("?q=");
    if (splitURL.length == 2 && splitURL[1] == '') {
        history.replaceState(null, null, splitURL[0])
    }
}

function toggle_serie_switch() {
    let serie_switch = <HTMLInputElement>document.getElementById("toggle_serie_mode");
    let url = window.location.href
    serie_switch.setAttribute("hx-get", url)
    // @ts-ignore
    htmx.process(serie_switch);
    // @ts-ignore
    htmx.trigger("#toggle_serie_mode","toggle_serie_switch");
}

document.addEventListener("htmx:configRequest", function(configEvent){
    let serie_switch = <HTMLInputElement>document.getElementById("toggle_serie_mode");
    // @ts-ignore
    configEvent.detail.headers['Smode'] = serie_switch.checked
})



