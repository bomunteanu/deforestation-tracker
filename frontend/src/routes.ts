import AreaDetails from './components/AreaDetails.svelte';
import AreaList from './components/AreaList.svelte';
import HistoryDetails from './components/HistoryDetails.svelte';
import Login from './components/Login.svelte'; // Import the Login component

export const routes = [
    { path: '/', component: AreaList },
    { path: '/areas/:id', component: AreaDetails },
    { path: '/history/:id', component: HistoryDetails },
    { path: '/login', component: Login } // Add the login route
];
