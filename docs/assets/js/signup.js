import { API_DOMAIN } from './config.js'

document.getElementById('signup-form').addEventListener('submit', async (e) => {
  e.preventDefault(); // Prevent the default form submission

  const name = document.getElementById('name').value;
  const password = document.getElementById('password').value;

  // Hide any previous error message
  const errorMessage = document.getElementById('error-message');
  const errorText = document.getElementById('error-text');
  errorMessage.classList.add('hidden');

  try {
    // First, signup the user
    const signupResponse = await fetch(`${API_DOMAIN}/sign-up`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, password }),
    });

    if (signupResponse.ok) {
      console.log("Signup successful, proceeding to login...");
        // Ensure login is awaited after signup success
      const loginResponse = await loginUser(name, password); // Await login function
      window.location.href = 'quote.html'
      if (loginResponse) {
        console.log("Login was successful, redirecting...");
        // Show success message and redirect
        showSuccessMessage('Signup and Login successful!');
        setTimeout(() => {
          window.location.href = 'quotes.html'; // Redirect after showing the message
        }, 2000); // Wait 2 seconds before redirecting
      } else {
        showErrorMessage('Login failed after signup. Please try again.');
      }
    } else {
      showErrorMessage('Signup failed. Please try again.');
    }
  } catch (error) {
    console.error('Error during signup or login:', error);
    showErrorMessage('An error occurred during signup or login. Please try again.');
  }
});

async function loginUser(name, password) {
  console.log("Attempting to login...");
  try {
    const loginResponse = await fetch(`${API_DOMAIN}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Ensure the cookie is included
      body: JSON.stringify({ name, password }),
    });

    if (loginResponse.ok) {
      const data = await loginResponse.json();
      console.log("Login response data: ", data); // Log the response data
      return data; // Return data to confirm login was successful
    } else {
      showErrorMessage("Login failed: " + await loginResponse.text());
      return null; // Return null if login failed
    }
  } catch (error) {
    console.error('Error during login:', error);
    showErrorMessage('Error during login: ' + error.message);
    return null; // Return null if an error occurred
  }
}

function showErrorMessage(message) {
  const errorMessage = document.getElementById('error-message');
  const errorText = document.getElementById('error-text');
  errorText.textContent = message;
  errorMessage.classList.remove('hidden');
}

function showSuccessMessage(message) {
  const successMessage = document.getElementById('success-message');
  const successText = document.getElementById('success-text');
  successText.textContent = message;
  successMessage.classList.remove('hidden');
}
