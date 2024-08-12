<script lang="ts">
    import { Pie } from 'svelte-chartjs';
    import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
    import type { ChartData } from 'chart.js';
  
    // Register Chart.js components
    ChartJS.register(ArcElement, Tooltip, Legend);
  
    export let deforestedArea: number;
    export let forestedArea: number;
  
    const data: ChartData<'pie'> = {
      labels: ['Deforested Area', 'Forested Area'],
      datasets: [
        {
          data: [deforestedArea, forestedArea],
          backgroundColor: ['#e74c3c', '#2ecc71'],
          borderColor: '#fff',
          borderWidth: 2,
        }
      ]
    };
  </script>
  
  <Pie {data} options={{ responsive: true, plugins: { legend: { position: 'top' } } }} />
  