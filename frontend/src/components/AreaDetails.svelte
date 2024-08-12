<script lang="ts">
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import APIClient from '../../sdk'; // Adjust the path if necessary
  import type { Area, History } from '../../sdk';
  import { calculateRectangleArea } from '../lib/area';
  import BarChart from './BarChart.svelte'; // Import the BarChart component

  const client = new APIClient('http://localhost:8080'); // Replace with your API base URL

  let area: Area | null = null;
  let histories: Array<{history: History, imageUrl: string | null}> = [];
  let loading = true;
  let error: string | null = null;

  // Extract the area ID from the route parameters
  let id: number;

  $: {
    // Extract the route parameter `id` from the URL
    id = Number(window.location.pathname.split('/').pop());
  }

  onMount(async () => {
    try {
      if (id) {
        area = await client.getArea(id);
        const fetchedHistories = await client.getHistoryByAreaId(id); // Fetch histories
        console.log(fetchedHistories);

        // Fetch images for each history
        for (const history of fetchedHistories) {
          try {
            const imageBlob = await client.getImageByPath(history.ImagePath.replace("/app/images/", ""));
            const imageUrl = URL.createObjectURL(imageBlob);
            console.log(imageUrl);
            histories.push({history, imageUrl});
          } catch (err) {
            histories.push({history, imageUrl: null});
            console.error('Failed to fetch image', err);
          }
        }
      }
    } catch (err) {
      error = 'Failed to fetch area details. Please try again later.';
      console.error(err);
    } finally {
      loading = false;
    }
  });

  const deleteArea = async () => {
    try {
      await client.deleteArea(id);
      navigate('/areas'); // Redirect to the areas list after deletion
    } catch (err) {
      error = 'Failed to delete area. Please try again later.';
      console.error(err);
    }
  };

  const goToHistoryDetails = (historyId: number) => {
    navigate(`/history/${historyId}`);
  };

  const goBack = () => {
    window.history.back(); // Use the browser's history API to go back
  };
</script>

<style>
  .area-details {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
    padding: 2rem;
    background-color: #ffffff;
    border-radius: 0.5rem;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
  }

  .text-section {
    flex: 1;
    max-width: 60%;
  }

  .chart-section {
    flex: 1;
    max-width: 40%;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: #f9f9f9;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 1rem;
  }

  .loading {
    font-style: italic;
    color: #666;
  }

  .error {
    color: #e74c3c;
    font-weight: bold;
  }

  .delete-button {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 0.75rem 1.5rem;
    background-color: #e74c3c;
    color: white;
    border: none;
    border-radius: 0.5rem;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.3s, transform 0.3s;
    z-index: 10;
  }

  .delete-button:hover {
    background-color: #c0392b;
    transform: scale(1.05);
  }

  .back-button {
    position: fixed;
    bottom: 1rem;
    left: 1rem;
    padding: 0.75rem 1.5rem;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 0.5rem;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.3s, transform 0.3s;
    z-index: 10;
  }

  .back-button:hover {
    background-color: #2980b9;
    transform: scale(1.05);
  }

  .history-list {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }

  .history-item {
    display: flex;
    align-items: center;
    padding: 1rem;
    border: 1px solid #ddd;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
    background-color: #f9f9f9;
    cursor: pointer;
    transition: background-color 0.3s, transform 0.3s;
  }

  .history-item:hover {
    background-color: #e9e9e9;
    transform: scale(1.02);
  }

  .text-content {
    flex: 1;
  }

  .image-content {
    margin-left: 1rem;
  }

  .history-image {
    max-width: 10rem;
    height: auto;
    border-radius: 0.5rem;
  }

  @media (max-width: 768px) {
    .area-details {
      flex-direction: column;
    }

    .text-section, .chart-section {
      max-width: 100%;
    }
  }
</style>

{#if loading}
  <p class="loading">Loading area details...</p>
{:else if error}
  <p class="error">{error}</p>
{:else if area}
  <div class="area-details">
    <div class="text-section">
      <h1>{area.AreaName}</h1>
      <p>Top Right Latitude: {area.TopRightLat}</p>
      <p>Top Right Longitude: {area.TopRightLon}</p>
      <p>Bottom Left Latitude: {area.BottomLeftLat}</p>
      <p>Bottom Left Longitude: {area.BottomLeftLon}</p>
      <p>Forested Area: {area.DeforestedArea !== undefined ? (100 - area.DeforestedArea).toFixed(2) + '%' : 'N/A'}</p>
      <button class="delete-button" on:click={deleteArea}>Delete Area</button>
    </div>
    <div class="chart-section">
      <BarChart {histories} />
    </div>
  </div>

  <ul class="history-list" style="padding-top: 4rem;">
    <h2 style="padding-top: 4rem;">History</h2>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    {#each histories as {history, imageUrl}}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
      <li class="history-item" on:click={() => goToHistoryDetails(history.ID)}>
        <div class="text-content">
          <p><strong>Date:</strong> {history.Date}</p>
          <p><strong>Forested Area:</strong> {((100 - history.DeforestedArea) * calculateRectangleArea(area.TopRightLat, area.TopRightLon, area.BottomLeftLat, area.BottomLeftLon)).toFixed(2)} km²</p>
        </div>
        <div class="image-content">
          {#if imageUrl}
            <!-- svelte-ignore a11y-img-redundant-alt -->
            <img src={imageUrl} alt="History Image" class="history-image" />
          {:else}
            <p>Image not available</p>
          {/if}
        </div>
      </li>
    {/each}
  </ul>

  <button class="back-button" on:click={goBack}>Back</button>
{/if}
