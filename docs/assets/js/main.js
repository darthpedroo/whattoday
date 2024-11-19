import { API_DOMAIN } from './config.js'

// Combine both navbar and quotes loading into one function
window.onload = async () => {
  await loadNavbar();  // Load the navbar first
  loadQuotes();        // Then load the quotes
};

// Function to load the navbar
async function loadNavbar() {
  const response = await fetch('navbar.html');
  const navbarHTML = await response.text();
  document.getElementById('navbar-container').innerHTML = navbarHTML;
}

// Select the container for quotes
const quoteContainer = document.getElementById('quote-container');

// Function to fetch and display quotes
async function loadQuotes() {
  try {
    const quotes = await fetchQuotes();
    if (quotes) {
      displayQuotes(quotes);
    } else {
      alert('Failed to load quotes');
    }
  } catch (error) {
    alert('An error occurred while fetching the quotes.');
    console.error('Error:', error);
  }
}

// Fetch the quotes from the API
async function fetchQuotes() {
  try {
    const response = await fetch(`${API_DOMAIN}/quotes`);
    if (response.ok) {
      return await response.json();
    }
    return null;
  } catch (error) {
    console.error('Error fetching quotes:', error);
    return null;
  }
}

// Display the quotes on the page
// Display the quotes on the page
function displayQuotes(quotes) {
  quoteContainer.innerHTML = '';
  quotes.forEach(quote => {
    const quoteDiv = document.createElement('div');
    quoteDiv.classList.add('quote', 'p-6', 'bg-white', 'border', 'border-gray-300', 'rounded-lg', 'mb-4', 'shadow-md');
    
    // Assuming quote has the fields: User.Name, Quote.Text, and PublishDate
    const publishDate = new Date(quote.Quote.PublishDate);
    const formattedDate = publishDate.toLocaleDateString() + ' ' + publishDate.toLocaleTimeString();

    quoteDiv.innerHTML = `
      <p><strong class="text-lg text-indigo-600">${quote.User.Name}</strong></p>
      <p class="text-gray-700 text-sm italic">— ${quote.Quote.Text}</p>
      <p class="text-xs text-gray-500 mt-2">Published on: ${formattedDate}</p>
    `;
    
    quoteContainer.appendChild(quoteDiv);
  });
}

