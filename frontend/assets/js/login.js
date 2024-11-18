document.getElementById('login-form').addEventListener('submit', async (e) => {
  e.preventDefault();

  const name = document.getElementById('name').value;
  const password = document.getElementById('password').value;

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
        console.log("EXXXXXXXXXXXXD")
        console.log(response)
      alert('Login successful!');
      window.location.href = 'quotes.html'; // Redirect to another page
    } else {
      alert('Invalid credentials. Please try again.');
    }
  } catch (error) {
    console.error('Error:', error);
    alert('An error occurred. Please try again.');
  }
});
