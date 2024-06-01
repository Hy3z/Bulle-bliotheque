enum SearchBoxState {
    Ilde,
    Writing,
}
let search_box_state = SearchBoxState.Ilde;
function update_search_history() {
    let search_box = <HTMLInputElement>document.getElementById("search-box");
    if (search_box_state == SearchBoxState.Ilde && search_box.value.length > 0) {
        history.pushState({}, "", "http://localhost:42069/browse?q="+search_box.value);
        search_box_state = SearchBoxState.Writing;

    }else if(search_box_state == SearchBoxState.Writing && search_box.value.length > 0) {
        history.replaceState({}, "", "http://localhost:42069/browse?q="+search_box.value);

    }else if(search_box_state == SearchBoxState.Writing && search_box.value.length == 0){
        history.pushState({}, "", "http://localhost:42069/browse");
        search_box_state = SearchBoxState.Ilde;
    }
}

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
}

window.addEventListener('unload', function(e){
    console.log("jhggfgbh")
});