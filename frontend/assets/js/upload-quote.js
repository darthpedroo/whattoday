document.getElementById('post-form').addEventListener('submit', async (e) => {
    e.preventDefault();
  
    const content = document.getElementById('content').value;
  
    try {
      const response = await fetch('http://localhost:8080/quotes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Include cookies in the request
        body: JSON.stringify({ text: content }), // Change 'content' to 'text'
      });
  
      if (response.ok) {
        const data = await response.json();
        alert('Post created successfully!');
        console.log('Created Post:', data);
        window.location.href = 'quotes.html'; // Redirect to a page showing all posts (placeholder)
      } else {
        alert('Failed to create the post. Please try again.');
      }
    } catch (error) {
      console.error('Error:', error);
      alert('An error occurred while creating the post.');
    }
});
