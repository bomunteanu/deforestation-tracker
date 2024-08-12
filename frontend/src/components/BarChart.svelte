<script lang="ts">
    import { Chart, registerables } from 'chart.js';
    import { Bar } from 'svelte-chartjs';
    import type { History } from '../../sdk'; // Adjust the path if necessary
  
    Chart.register(...registerables);
  
    export let histories: Array<{ history: History }> = [];
  
    // Prepare chart data and options
    let data = {
      labels: histories.map(h => h.history.Date),
      datasets: [{
        label: 'Forestation (%)',
        data: histories.map(h => 100 - h.history.DeforestedArea),
        backgroundColor: 'rgba(75, 192, 192, 0.2)', // Light green
        borderColor: 'rgba(75, 192, 192, 1)', // Dark green
        borderWidth: 1
      }]
    };
  
    // Type for Chart.js v3 or v4
    type TooltipItem = import('chart.js').TooltipItem<'bar'>;
    type TooltipModel = import('chart.js').TooltipModel<'bar'>;
  
    let options = {
      responsive: true,
      plugins: {
        legend: {
          position: 'top' as const,
          labels: {
            font: {
              size: 14,
              weight: 'bold'
            }
          }
        },
        tooltip: {
          callbacks: {
            label: function(tooltipItem: TooltipItem<'bar'>) {
              const label = tooltipItem.dataset.label || 'No label';
              const value = tooltipItem.raw.toFixed(2);
              return `${label}: ${value} %`;
            }
          },
          bodyFont: {
            size: 12
          },
          titleFont: {
            size: 14,
            weight: 'bold'
          }
        }
      },
      scales: {
        x: {
          beginAtZero: true,
          grid: {
            color: '#e0e0e0', // Light gray grid lines
            borderColor: '#d0d0d0' // Light border color for x-axis
          },
          ticks: {
            font: {
              size: 12
            }
          }
        },
        y: {
          beginAtZero: true,
          min: 0,
          max: 100, // Set the maximum value of the y-axis to 100
          title: {
            display: true,
            text: 'Forestation (%)',
            font: {
              size: 14,
              weight: 'bold'
            }
          },
          grid: {
            color: '#e0e0e0', // Light gray grid lines
            borderColor: '#d0d0d0' // Light border color for y-axis
          },
          ticks: {
            font: {
              size: 12
            }
          }
        }
      }
    };
  </script>
  
  <style>
    .chart-container {
      width: 100%;
      max-width: 800px;
      margin: 0 auto;
      padding: 1rem; /* Add padding around the chart */
      background: #f9f9f9; /* Light background color */
      border-radius: 8px; /* Rounded corners */
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* Subtle shadow */
    }
  
    canvas {
      width: 100%;
      height: 100px; /* Fixed height */
    }
  </style>
  
  <div class="chart-container">
    <!-- Pass data and options as props to the Bar component -->
    <Bar {data} {options} />
  </div>
  