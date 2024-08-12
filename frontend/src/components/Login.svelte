<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import APIClient from '../../sdk'; // Adjust the path if necessary

    const dispatch = createEventDispatcher();
    let email = '';
    let password = '';
    let error: string | null = null;
    let loading = false;

    const handleLogin = async () => {
      loading = true;
      error = null;

      try {
        const client = new APIClient('http://localhost:8080'); // Adjust base URL if needed
        await client.login(email, password);
        dispatch('loginSuccess');
        window.location.replace("/");
      } catch (err) {
        error = 'Invalid username or password. Please try again.';
        console.error(err);
      } finally {
        loading = false;
      }
    };
</script>

<style>
    .login-container {
      max-width: 400px;
      margin: 0 auto;
      padding: 2rem;
      background-color: #ffffff;
      border-radius: 0.5rem;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .login-container h1 {
      font-size: 1.5rem;
      margin-bottom: 1rem;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .form-group label {
      font-weight: bold;
    }

    .form-group input {
      padding: 0.75rem;
      border: 1px solid #ddd;
      border-radius: 0.25rem;
      font-size: 1rem;
    }

    .error-message {
      color: #e74c3c;
      font-weight: bold;
    }

    .submit-button {
      padding: 0.75rem;
      background-color: #3498db;
      color: white;
      border: none;
      border-radius: 0.25rem;
      cursor: pointer;
      font-size: 1rem;
      transition: background-color 0.3s, transform 0.3s;
    }

    .submit-button:hover {
      background-color: #2980b9;
      transform: scale(1.05);
    }

    .submit-button:disabled {
      background-color: #95a5a6;
      cursor: not-allowed;
    }

    .loading {
      font-style: italic;
      color: #666;
    }
</style>

<div class="login-container">
    <h1>Login</h1>
    {#if error}
      <p class="error-message">{error}</p>
    {/if}
    {#if loading}
      <p class="loading">Logging in...</p>
    {/if}
    <div class="form-group">
      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} placeholder="Enter your email" />
    </div>
    <div class="form-group">
      <label for="password">Password</label>
      <input id="password" type="password" bind:value={password} placeholder="Enter your password" />
    </div>
    <button class="submit-button" on:click={handleLogin} disabled={loading}>Log In</button>
</div>
