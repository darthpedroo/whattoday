import { loginUser } from './auth.js';

document.getElementById('login-form').addEventListener('submit', async (e) => {
  e.preventDefault();

  const name = document.getElementById('name').value;
  const password = document.getElementById('password').value;

  const loginResponse = await loginUser(name, password);

  if (loginResponse) {
    console.log("Login successful, redirecting...");
    window.location.href = 'quote.html';
  } else {
    const errorMessage = document.getElementById('error-message');
    const errorText = document.getElementById('error-text');
    errorText.textContent = 'Login failed. Please try again.';
    errorMessage.classList.remove('hidden');
  }
});
