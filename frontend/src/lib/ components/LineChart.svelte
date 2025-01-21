<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, Tooltip, Legend } from 'chart.js';

    // Register necessary Chart.js components
    Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, Tooltip, Legend);

    export let stats : DataSet[] 
    export let label : string

    let chart: Chart | null = null
    Chart.defaults.color = 'rgb(250,255,255)'
    Chart.defaults.font.size = 14
    let chartCanvas: HTMLCanvasElement
    $: {
        if (chart) {
            chart.destroy()
        }

        if (stats.length > 0) {
            chart = new Chart(chartCanvas, {
                type: 'line',
                data: {
                    labels: stats[0].xAxis,
                    datasets: stats.map(stat => ({
                        label: stat.label,
                        data: stat.data,
                        borderColor: stat.borderColor || 'rgba(75, 192, 192, 1)',
                        backgroundColor: stat.backgroundColor || 'rgba(75, 192, 192, 0.2)',
                        fill: true,
                    }))
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            title: {
                                display: true,
                                text: label,
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: 'Date'
                            }
                        }
                    },
                }
            })
        }
    }
    
    onDestroy(() => {
        if (chart) chart.destroy()
    })
</script>

<canvas class="w-full" bind:this={chartCanvas}></canvas>