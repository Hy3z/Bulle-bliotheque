// Fonction pour basculer la visibilité du menu déroulant
function toggleDropdown(): void {
  const dropdown = document.getElementById("myDropdown");
  // Vérifiez si l'élément existe avant d'appeler classList
  if (dropdown) {
    dropdown.classList.toggle("show");
  }
}

window.onclick = function(event) {
  console.log("test");
  if (!event.target.closest('.dropbtn, .dropdown-content')) {
    const dropdowns = document.getElementsByClassName("dropdown-content");
    for (let i = 0; i < dropdowns.length; i++) {
      const openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
        openDropdown.classList.remove('show');
      }
    }
  }
};


