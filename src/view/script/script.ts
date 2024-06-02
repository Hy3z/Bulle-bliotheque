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
/*
enum BrowseBarState {
    Ilde,
    Writing,
}
let search_box_state = BrowseBarState.Ilde;
function update_search_history() {
    let browse_bar = <HTMLInputElement>document.getElementById("browse-bar");
    if (search_box_state == BrowseBarState.Ilde && browse_bar.value.length > 0) {
        history.pushState({}, "", "http://localhost:42069/browse?q="+browse_bar.value);
        search_box_state = BrowseBarState.Writing;

    }else if(search_box_state == BrowseBarState.Writing && browse_bar.value.length > 0) {
        history.replaceState({}, "", "http://localhost:42069/browse?q="+browse_bar.value);

    }else if(search_box_state == BrowseBarState.Writing && browse_bar.value.length == 0){
        history.pushState({}, "", "http://localhost:42069/browse");
        search_box_state = BrowseBarState.Ilde;
    }
}*/

//SerieSwitch ---------------------
/*
function toggle_serie_switch() {
    let serie_switch = <HTMLInputElement>document.getElementById("toggle_serie_mode");
    let series = document.getElementsByClassName("serie-preview");
    let books_with_serie = document.getElementsByClassName("hidable-book-preview");
    for (let i = 0; i < series.length; i++) {
        (series[i] as HTMLElement).style.display = serie_switch.checked ? 'block' : 'none';
    }
    for (let i = 0; i < books_with_serie.length; i++) {
        (books_with_serie[i] as HTMLElement).style.display = serie_switch.checked ? 'none' : 'block';
    }
}*/
function toggle_serie_switch() {
    let serie_switch = <HTMLInputElement>document.getElementById("toggle_serie_mode");
    let url = window.location.href
    serie_switch.setAttribute("hx-get", url)
    serie_switch.setAttribute("hx-headers", `{"Tmpl":"main", "SerieMode":"${serie_switch.checked}"`)
    console.log("EYO")
    // @ts-ignore
    htmx.process(serie_switch)
    // @ts-ignore
    htmx.trigger("#toggle_serie_mode","toggle_serie_switch")
}



