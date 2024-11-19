document.getElementById('login-form').addEventListener('submit', async (e) => {
  e.preventDefault();

  const name = document.getElementById('name').value;
  const password = document.getElementById('password').value;

  // Hide any previous error message
  const errorMessage = document.getElementById('error-message');
  const errorText = document.getElementById('error-text');
  errorMessage.classList.add('hidden');

  try {
    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies in the request
      body: JSON.stringify({ name, password }),
    });

    if (response.ok) {
      // Redirect on successful login
      window.location.href = 'quote.html'; // Redirect to another page
    } else {
      // Show error message if login fails
      errorText.textContent = 'Invalid credentials. Please try again.';
      errorMessage.classList.remove('hidden');
    }
  } catch (error) {
    console.error('Error:', error);
    errorText.textContent = 'An error occurred. Please try again.';
    errorMessage.classList.remove('hidden');
  }
});
