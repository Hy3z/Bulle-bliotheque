//Browse bar ---------------------
var BrowseBarState;
(function (BrowseBarState) {
    BrowseBarState[BrowseBarState["Ilde"] = 0] = "Ilde";
    BrowseBarState[BrowseBarState["Writing"] = 1] = "Writing";
})(BrowseBarState || (BrowseBarState = {}));
var search_box_state = BrowseBarState.Ilde;
function update_search_history() {
    var browse_bar = document.getElementById("browse-bar");
    if (search_box_state == BrowseBarState.Ilde && browse_bar.value.length > 0) {
        history.pushState({}, "", "http://localhost:42069/browse?q=" + browse_bar.value);
        search_box_state = BrowseBarState.Writing;
    }
    else if (search_box_state == BrowseBarState.Writing && browse_bar.value.length > 0) {
        history.replaceState({}, "", "http://localhost:42069/browse?q=" + browse_bar.value);
    }
    else if (search_box_state == BrowseBarState.Writing && browse_bar.value.length == 0) {
        history.pushState({}, "", "http://localhost:42069/browse");
        search_box_state = BrowseBarState.Ilde;
    }
}
function clear_browse_bar() {
    var browse_bar = document.getElementById("browse-bar");
    browse_bar.value = '';
}
//SerieSwitch ---------------------
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
