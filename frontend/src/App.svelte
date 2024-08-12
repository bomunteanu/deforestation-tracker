<script lang="ts">
  import { Router, Route, Link } from 'svelte-routing';
  import { routes } from './routes';
  import RouteGuard from './components/RouteGuard.svelte';
  import { isAuthenticated } from './store'; // Adjust the path if necessary


  // This function checks if a route is protected
  const isProtectedRoute = (path: string) => {
    return !['/login'].includes(path); // Add other public routes if necessary
  };

  // Function to handle logout
  const handleLogout = () => {
    localStorage.removeItem('authToken'); // Clear the token from local storage
    isAuthenticated.set(false); // Update the authentication state
    window.location.replace("/login"); // Redirect to login page
  };
</script>

<style>
  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 2rem;
    background-color: #0a7e2f;
    color: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    position: sticky;
    top: 0;
    z-index: 1000;
  }

  .nav-brand {
    font-size: 1.5rem;
    font-weight: bold;
  }

  .nav-links {
    display: flex;
    gap: 1.5rem;
  }

  .nav-links a {
    color: #fff;
    text-decoration: none;
    font-size: 1rem;
    transition: color 0.3s;
  }

  .nav-links a:hover {
    color: #f0f0f0;
  }

  .nav-links a.active {
    border-bottom: 2px solid #f0f0f0;
    color: #f0f0f0;
    padding-bottom: 0.2rem;
  }

  main {
    padding: 2rem;
  }
</style>

<Router>
  <nav>
    <div class="nav-brand">Doftana Deforestation Tracker</div>
    <div class="nav-links">
      <Link to="/" exact activeClassName="active">Home</Link>
      {#if $isAuthenticated}
        <a href="#" on:click={handleLogout}>Logout</a>
      {:else}
        <Link to="/login" activeClassName="active">Login</Link>
      {/if}
    </div>
  </nav>
  <main>
    {#each routes as { path, component }}
      {#if isProtectedRoute(path)}
        <RouteGuard>
          <Route path={path} component={component} />
        </RouteGuard>
      {:else}
        <Route path={path} component={component} />
      {/if}
    {/each}
  </main>
</Router>
