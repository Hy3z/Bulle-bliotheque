var SearchBoxState;
(function (SearchBoxState) {
    SearchBoxState[SearchBoxState["Ilde"] = 0] = "Ilde";
    SearchBoxState[SearchBoxState["Writing"] = 1] = "Writing";
})(SearchBoxState || (SearchBoxState = {}));
var search_box_state = SearchBoxState.Ilde;
function update_search_history() {
    var search_box = document.getElementById("search-box");
    if (search_box_state == SearchBoxState.Ilde && search_box.value.length > 0) {
        history.pushState({}, "", "http://localhost:42069/browse?q=" + search_box.value);
        search_box_state = SearchBoxState.Writing;
    }
    else if (search_box_state == SearchBoxState.Writing && search_box.value.length > 0) {
        history.replaceState({}, "", "http://localhost:42069/browse?q=" + search_box.value);
    }
    else if (search_box_state == SearchBoxState.Writing && search_box.value.length == 0) {
        history.pushState({}, "", "http://localhost:42069/browse");
        search_box_state = SearchBoxState.Ilde;
    }
}
function toggle_serie_switch() {
    var serie_switch = document.getElementById("toggle_serie_mode");
    var series = document.getElementsByClassName("serie-preview");
    var books_with_serie = document.getElementsByClassName("hidable-book-preview");
    for (var i = 0; i < series.length; i++) {
        series[i].style.display = serie_switch.checked ? 'block' : 'none';
    }
    for (var i = 0; i < books_with_serie.length; i++) {
        books_with_serie[i].style.display = serie_switch.checked ? 'none' : 'block';
    }
}
window.addEventListener('unload', function (e) {
    console.log("jhggfgbh");
});
