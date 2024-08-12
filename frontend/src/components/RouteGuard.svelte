<script lang="ts">
    import { onMount } from 'svelte';
    import { navigate } from 'svelte-routing';
    import APIClient from '../../sdk'; // Adjust the path if necessary
    import { isAuthenticated } from '../store'; // Adjust the path if necessary
    import { get } from 'svelte/store';
  
    const client = new APIClient('http://localhost:8080'); // Adjust base URL if needed
  
    // Function to check if the user is authenticated
    const checkAuthentication = async () => {
      try {
        await client.check();
        isAuthenticated.set(true);
      } catch (err) {
        isAuthenticated.set(false);
      }
    };
  
    // Redirect to login if not authenticated
    onMount(async () => {
      await checkAuthentication();
      if (!get(isAuthenticated)) {
        navigate('/login');
      }
    });
  </script>
  
  {#if $isAuthenticated}
    <slot />
  {/if}
  