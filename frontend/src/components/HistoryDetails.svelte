<script lang="ts">
    import { onMount } from 'svelte';
    import { navigate } from 'svelte-routing';
    import APIClient from '../../sdk'; // Adjust the path if necessary
    import type { Area, History } from '../../sdk';
    import { calculateRectangleArea } from '../lib/area';
    import PieChart from './PieChart.svelte'; // Import the PieChart component
  
    const client = new APIClient('http://34.88.213.170:8080'); // Replace with your API base URL
  
    let history: History | null = null;
    let area: Area | null = null;
    let imageUrl: string | null = null;
    let maskedImageUrl: string | null = null;
    let loading = true;
    let error: string | null = null;
  
    // Extract the history ID from the route parameters
    let id: number;
  
    $: {
      // Extract the route parameter `id` from the URL
      id = Number(window.location.pathname.split('/').pop());
    }
  
    onMount(async () => {
      try {
        if (id) {
          history = await client.getHistoryById(id);
          area = await client.getArea(history.AreaID);
          console.log(history);
          try {
            const imageBlob = await client.getImageByPath(history.ImagePath.replace("/app/images/", ""));
            const maskedImageBlob = await client.getImageByPath(history.MaskedImagePath.replace("/app/images/", ""));
            imageUrl = URL.createObjectURL(imageBlob);
            maskedImageUrl = URL.createObjectURL(maskedImageBlob);
            console.log(imageUrl);
          } catch (err) {
            console.error('Failed to fetch image', err);
          }
        }
      } catch (err) {
        error = 'Failed to fetch history details. Please try again later.';
        console.error(err);
      } finally {
        loading = false;
      }
    });
  
    const goBack = () => {
      window.history.back(); // Use the browser's history API to go back
    };
  
    $: deforestedPercentage = history ? history.DeforestedArea : 0;
    $: totalArea = area ? calculateRectangleArea(area.TopRightLat, area.TopRightLon, area.BottomLeftLat, area.BottomLeftLon) : 0;
    $: forestedArea = history ? (100 - deforestedPercentage) * totalArea / 100 : 0;
  </script>
  
  <style>
    .history-details {
      display: flex;
      gap: 1rem;
      padding: 1rem;
      background-color: #f9f9f9; /* Light background for better contrast */
      border-radius: 0.5rem; /* Rounded corners for the details container */
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* Subtle shadow for depth */
    }
  
    .details-left {
      flex: 1;
      max-width: 400px;
    }
  
    .details-right {
      flex: 2;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
  
    .loading {
      font-style: italic;
      color: #555; /* Slightly darker color for better readability */
    }
  
    .error {
      color: #e74c3c; /* Red color for error messages */
      font-weight: bold;
    }
  
    .back-button {
      margin-bottom: 1rem;
      padding: 0.5rem 1rem;
      background-color: #007bff;
      color: white;
      border: none;
      border-radius: 0.25rem;
      cursor: pointer;
      font-size: 1rem;
      text-align: center;
    }
  
    .back-button:hover {
      background-color: #0056b3;
    }
  
    .history-details h1 {
      margin-top: 0;
      font-size: 1.5rem;
      color: #333; /* Darker color for the heading */
    }
  
    .history-details p {
      margin: 0.5rem 0; /* Space between paragraphs */
      font-size: 1rem;
      color: #666; /* Medium gray for text */
    }
  
    .history-image,
    .masked-image {
      width: 100%; /* Ensure images take up the full width of their container */
      max-width: 600px; /* Limit the maximum width */
      height: auto; /* Maintain aspect ratio */
      border-radius: 0.5rem; /* Rounded corners for images */
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* Subtle shadow for depth */
    }
  
    .history-image {
      margin-bottom: 1rem; /* Space between images */
    }
  
    .masked-image {
      margin-bottom: 1rem; /* Space below the masked image */
    }
  
    .chart-container {
      margin-bottom: 1rem; /* Padding between chart and details */
    }
  </style>
  
  {#if loading}
    <p class="loading">Loading history details...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else if history}
    <div class="history-details">
      <div class="details-left" style="padding-right: 5rem;">
        <button class="back-button" on:click={goBack}>Back</button>
        <div style="margin-bottom: 5rem;">
            <h1>History Details</h1>
            <p><strong>Date:</strong> {history.Date}</p>
            <p><strong>Forested Area:</strong> {forestedArea.toFixed(2)} km²</p>
        </div>
        <div class="chart-container">
          <PieChart deforestedArea={deforestedPercentage * totalArea / 100} forestedArea={forestedArea} />
        </div>
      </div>
      <div class="details-right">
        {#if imageUrl}
          <!-- svelte-ignore a11y-img-redundant-alt -->
          <img src={imageUrl} alt="History Image" class="history-image" />
          <!-- svelte-ignore a11y-img-redundant-alt -->
          <img src={maskedImageUrl} alt="Masked Image" class="masked-image" />
        {:else}
          <p>Image not available</p>
        {/if}
      </div>
    </div>
  {/if}
  
