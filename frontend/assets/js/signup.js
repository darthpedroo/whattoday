document.getElementById('signup-form').addEventListener('submit', async (e) => {
    e.preventDefault();
  
    const name = document.getElementById('name').value;
    const password = document.getElementById('password').value;
  
    try {
      const signupResponse = await fetch('http://localhost:8080/sign-up', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, password }),
      });
  
      if (signupResponse.ok) {
        // Automatically login after signup
        const loginResponse = await fetch('http://localhost:8080/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ name, password }),
        });
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
              body: JSON.stringify({ name, password }),
            });
        
            if (response.ok) {
              const data = await response.json();
              alert('Login successful!');
              console.log('Token:', data.token); // Debugging only
              window.location.href = 'quotes.html'; // Redirect to another page
            } else {
              alert('Invalid credentials. Please try again.');
            }
          } catch (error) {
            console.error('Error:', error);
            alert('An error occurred. Please try again.');
          }
        });
        
        if (loginResponse.ok) {
          const data = await loginResponse.json();
          alert('Signup and Login successful!');
          console.log('Token:', data.token); // Debugging only
          window.location.href = 'quotes.html'; // Redirect to another page
        } else {
          alert('Signup succeeded, but login failed.');
        }
      } else {
        alert('Signup failed. Please try again.');
      }
    } catch (error) {
      console.error('Error:', error);
      alert('An error occurred. Please try again.');
    }
  });
  