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
    const response = await fetch('http://localhost:8080/quotes');
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
function displayQuotes(quotes) {
  quoteContainer.innerHTML = '';
  quotes.forEach(quote => {
    const quoteDiv = document.createElement('div');
    quoteDiv.classList.add('quote');
    quoteDiv.innerHTML = `
      <p><strong>${quote.User.Name}</strong></p>
      <p><em>— ${quote.Quote.Text}</em></p>
    `;
    quoteContainer.appendChild(quoteDiv);
  });
}

// Call the loadQuotes function when the page loads
window.onload = loadQuotes;
