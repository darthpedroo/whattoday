import { API_DOMAIN } from './config.js'

document.addEventListener('DOMContentLoaded', async () => {
  // Check if the user is logged in by verifying the JWT from cookies
  const isLoggedIn = await checkLoginStatus();

  if (isLoggedIn) {
    // Show the post creation form if logged in
    document.getElementById('post-form').classList.remove('hidden');
  } else {
    // Show the login message if not logged in
    document.getElementById('login-message').classList.remove('hidden');
  }
});

document.getElementById('post-form').addEventListener('submit', async (e) => {
  e.preventDefault();

  const content = document.getElementById('content').value;

  // Hide any previous messages
  const errorMessage = document.getElementById('error-message');
  const successMessage = document.getElementById('success-message');
  const errorText = document.getElementById('error-text');
  errorMessage.classList.add('hidden');
  successMessage.classList.add('hidden');

  try {
    const response = await fetch(`${API_DOMAIN}/quotes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify({ text: content }), // Change 'content' to 'text'
    });

    if (response.ok) {
      const data = await response.json();
      console.log('Created Post:', data);
      window.location.href = 'quotes.html'

      // Show success message
      successMessage.classList.remove('hidden');
      setTimeout(() => {
        window.location.href = 'quotes.html'; // Redirect to a page showing all posts (placeholder)
      }, 2000); // Wait 2 seconds before redirecting
    } else {
      // Attempt to parse the error message from the server response
      const errorResponse = await response.json();
      const errorMessageText = errorResponse.error || 'Failed to create the post. Please try again.';
      
      // Show error message from server response
      errorText.textContent = errorMessageText;
      errorMessage.classList.remove('hidden');
    }
  } catch (error) {
    console.error('Error:', error);

    // Show generic error message if there's a network or other unexpected error
    errorText.textContent = 'An error occurred while creating the post.';
    errorMessage.classList.remove('hidden');
  }
});

// Function to check if the user is logged in by validating the JWT
async function checkLoginStatus() {
  const jwt = getJWTFromCookies();

  if (!jwt) {
    return false; // No JWT, user is not logged in
  }

  const isValid = validateJWT(jwt);
  return isValid;
}

// Function to retrieve the JWT from cookies
function getJWTFromCookies() {
  // Get cookies as a string and split them into individual cookies
  const cookies = document.cookie.split(';');
  
  // Find the cookie named 'Authorization' and return its value
  for (let cookie of cookies) {
    cookie = cookie.trim(); // Remove leading spaces
    console.log(cookie)
    console.log("caca")
    if (cookie.startsWith('Authorization=')) {
      return cookie.substring('Authorization='.length); // Return the JWT token
    }
  }
  return null; // Return null if the 'Authorization' cookie is not found
}

// Function to validate JWT (check expiration)
function validateJWT(jwt) {
  const payload = parseJWT(jwt);
  if (!payload) return false;

  // Check if the token has expired
  const currentTime = Math.floor(Date.now() / 1000); // Current time in seconds
  return payload.exp > currentTime;
}

// Function to parse JWT and extract the payload
function parseJWT(jwt) {
  const base64Url = jwt.split('.')[1]; // The payload is the second part of the JWT
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Decode base64

  try {
    const jsonPayload = decodeURIComponent(escape(window.atob(base64))); // Decode and handle UTF-8
    return JSON.parse(jsonPayload);
  } catch (e) {
    console.error('Error parsing JWT:', e);
    return null;
  }
}
