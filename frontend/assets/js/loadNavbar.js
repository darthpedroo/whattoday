async function loadNavbar() {
    const response = await fetch('navbar.html');
    const navbarHTML = await response.text();
    document.getElementById('navbar-container').innerHTML = navbarHTML;
  }

  // Load the navbar when the page loads
  window.onload = loadNavbar;