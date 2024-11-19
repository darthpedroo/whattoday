// auth.js
import { API_DOMAIN } from './config.js';

export async function loginUser(name, password) {
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

      // Set the cookie manually in the frontend
      document.cookie = `Authorization=${data.token}; path=/; Secure; SameSite=None`;

      return data; // Return data to confirm login was successful
    } else {
      console.error("Login failed: ", await loginResponse.text());
      return null; // Return null if login failed
    }
  } catch (error) {
    console.error('Error during login:', error);
    return null; // Return null if an error occurred
  }
}
