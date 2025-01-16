<script lang="ts">
	import LineChart from "$lib/ components/LineChart.svelte";
	export let data: {
		records: townRecords[]
	}
</script>

<div class="flex flex-col justify-center items-center h-screen w-full">
	<LineChart
	{...{
		stats: [
			{
				label: 'Historical Bid-Ask Spread Percentage',
				data: [1, 2, 3],
				xAxis: [1, 2, 3],
				borderColor: 'rgba(75, 192, 192, 1)',
				backgroundColor: 'rgba(75, 192, 192, 0.2)',
			}
		],
		label: "Hi"
	}}>

	</LineChart>
</div>