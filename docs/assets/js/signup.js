import { API_DOMAIN } from './config.js';
import { loginUser } from './auth.js';

document.getElementById('signup-form').addEventListener('submit', async (e) => {
  e.preventDefault();

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

      // Use imported loginUser function
      const loginResponse = await loginUser(name, password);

      if (loginResponse) {
        console.log("Login was successful, redirecting...");
        window.location.href = 'quote.html'; // Redirect after successful login
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

function showErrorMessage(message) {
  const errorMessage = document.getElementById('error-message');
  const errorText = document.getElementById('error-text');
  errorText.textContent = message;
  errorMessage.classList.remove('hidden');
}
