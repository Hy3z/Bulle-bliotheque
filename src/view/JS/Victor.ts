// Fonction pour basculer la visibilité du menu déroulant
function toggleDropdown(): void {
    const dropdown = document.getElementById("myDropdown");
    // Vérifiez si l'élément existe avant d'appeler classList
    if (dropdown) {
      dropdown.classList.toggle("show");
    }
  }
  
  // Ferme le menu si l'utilisateur clique en dehors de celui-ci
  window.onclick = function(event: MouseEvent): void {
    // Utilisez l'opérateur de coalescence des nuls pour gérer les cas où event.target est null
    const target = event.target as Element;
  
    // Cette vérification s'assure que l'élément sur lequel l'utilisateur a cliqué n'est pas un bouton et n'est pas null
    if (!target.matches('img')) {
      const dropdowns = document.getElementsByClassName("dropdown-content");
      
      for (let i = 0; i < dropdowns.length; i++) {
        const openDropdown = dropdowns[i] as HTMLElement;
  
        if (openDropdown.classList.contains('show')) {
          openDropdown.classList.remove('show');
        }
      }
    }
  }
  
  