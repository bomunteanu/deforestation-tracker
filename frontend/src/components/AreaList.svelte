<!-- src/routes/Home.svelte -->
<script lang="ts">
    import { onMount, afterUpdate } from 'svelte';
    import { navigate } from 'svelte-routing'; // Import navigate for navigation
    import APIClient from '../../sdk'; // Adjust path if necessary
    import type { Area } from '../../sdk';
    import L from 'leaflet';
    import { calculateRectangleArea } from '../lib/area';

    const client = new APIClient('http://localhost:8080'); // Replace with your API base URL

    let areas: Array<Area> = [];
    let loading = true;
    let error: string | null = null;

    // Form visibility
    let showForm = false;

    // Form data
    let newArea = {
        area_name: '',
        top_right_lat: 45.379264,
        top_right_lon: 25.713041,
        bottom_left_lat: 45.369610,
        bottom_left_lon: 25.693347
    };

    // Map variables
    let mapContainer: HTMLDivElement | null = null;
    let map: L.Map;
    let topRightMarker: L.Marker;
    let bottomLeftMarker: L.Marker;

    onMount(async () => {
        try {
            // Fetch all areas
            areas = await client.getAllAreas();
            console.log(areas)
        } catch (err) {
            error = 'Failed to fetch areas. Please try again later.';
            console.error(err);
        } finally {
            loading = false;
        }
        
        // Initialize the map
        initializeMap();
    });

    afterUpdate(() => {
        if (mapContainer && !map) {
            initializeMap();
        }
    });

    const initializeMap = () => {
        if (!mapContainer) {
            console.error("Map container not found.");
            return;
        }

        map = L.map(mapContainer).setView([45.374437, 25.703194], 13);

        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);

        const redPinIcon = L.icon({
            iconUrl: 'src/assets/red_pin.png', // Replace with the URL to your pin icon image
            iconSize: [38, 38], // Size of the icon
            iconAnchor: [19, 38], // Point of the icon which will correspond to marker's location
            popupAnchor: [0, -38] // Point from which the popup should open relative to the iconAnchor
        });

        const bluePinIcon = L.icon({
            iconUrl: 'src/assets/blue_pin.png', // Replace with the URL to your pin icon image
            iconSize: [38, 38], // Size of the icon
            iconAnchor: [19, 38], // Point of the icon which will correspond to marker's location
            popupAnchor: [0, -38] // Point from which the popup should open relative to the iconAnchor
        });

        // Add default markers with updated positions
        topRightMarker = L.marker([newArea.top_right_lat, newArea.top_right_lon], { draggable: true, icon: redPinIcon }).addTo(map)
            .bindPopup('Top Right Corner')
            .on('dragend', () => updateCoordinates());

        bottomLeftMarker = L.marker([newArea.bottom_left_lat, newArea.bottom_left_lon], { draggable: true, icon: bluePinIcon }).addTo(map)
            .bindPopup('Bottom Left Corner')
            .on('dragend', () => updateCoordinates());
    };

    const updateCoordinates = () => {
        if (!topRightMarker || !bottomLeftMarker) return;

        const topRightLatLng = topRightMarker.getLatLng();
        const bottomLeftLatLng = bottomLeftMarker.getLatLng();

        newArea.top_right_lat = topRightLatLng.lat;
        newArea.top_right_lon = topRightLatLng.lng;
        newArea.bottom_left_lat = bottomLeftLatLng.lat;
        newArea.bottom_left_lon = bottomLeftLatLng.lng;
    };

    const createArea = async () => {
        if (!newArea.area_name.trim()) {
            alert('Area name cannot be empty.');
            return;
        }

        try {
            await client.createArea(newArea);
            // Refresh the area list
            areas = await client.getAllAreas();
            // Hide the form
            showForm = false;
            // Clear form data
            newArea = {
                area_name: '',
                top_right_lat: 45.379264,
                top_right_lon: 25.713041,
                bottom_left_lat: 45.369610,
                bottom_left_lon: 25.693347
            };
        } catch (err) {
            error = 'Failed to create area. Please try again later.';
            console.error(err);
        }
    };

    const viewAreaDetails = (areaId: string) => {
        navigate(`/areas/${areaId}`);
    };
</script>

<style>
    .area-list {
        list-style: none;
        padding: 0;
    }

    .area-item {
        margin-bottom: 1rem;
        padding: 1rem;
        border: 1px solid #ddd;
        border-radius: 0.5rem;
        cursor: pointer; /* Add a pointer cursor */
    }

    .loading {
        font-style: italic;
    }

    .error {
        color: red;
        font-weight: bold;
    }

    .form-container {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .form {
        background: white;
        padding: 2rem;
        border-radius: 0.5rem;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        width: 600px;
        max-width: 90%;
        position: relative;
    }

    .form button {
        margin-top: 1rem;
    }

    .create-button {
        margin: 1rem 0;
        padding: 0.5rem 1rem;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 0.25rem;
        cursor: pointer;
    }

    .create-button:hover {
        background-color: #0056b3;
    }

    #map {
        height: 300px;
        width: 100%;
        margin-bottom: 1rem;
    }
</style>

{#if loading}
    <p class="loading">Loading areas...</p>
{:else if error}
    <p class="error">{error}</p>
{:else}
    <button class="create-button" on:click={() => showForm = true}>Create New Area</button>
    {#if showForm}
        <div class="form-container">
            <div class="form">
                <h2>Create New Area</h2>
                <label>
                    Area Name:
                    <input type="text" bind:value={newArea.area_name} />
                </label>
                <p>Red pin is for top right, Blue pin is for bottom left</p>
                <div id="map" bind:this={mapContainer}></div>
                <button on:click={createArea}>Submit</button>
                <button on:click={() => showForm = false}>Cancel</button>
            </div>
        </div>
    {/if}
    <ul class="area-list">
        {#each areas as area}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
            <li class="area-item" on:click={() => viewAreaDetails(JSON.stringify(area.ID))}>
                <h3>{area.AreaName}</h3>
                <p>Forested Area: {((100 - area.DeforestedArea) * calculateRectangleArea(area.TopRightLat, area.TopRightLon, area.BottomLeftLat, area.BottomLeftLon)).toFixed(2)} km^2</p>
            </li>
        {/each}
    </ul>
{/if}
