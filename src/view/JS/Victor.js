// Fonction pour basculer la visibilité du menu déroulant
function toggleDropdown() {
    var dropdown = document.getElementById("myDropdown");
    // Vérifiez si l'élément existe avant d'appeler classList
    if (dropdown) {
        dropdown.classList.toggle("show");
    }
}
window.onclick = function (event) {
    console.log("test");
    if (!event.target.closest('.dropbtn, .dropdown-content')) {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        for (var i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
};
