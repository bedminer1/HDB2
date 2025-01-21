<script lang="ts">
	import LineChart from "$lib/ components/LineChart.svelte";
	export let data: {
		records: townRecords[]
		start: string
		end: string
	}

	let start = data.start
	let end = data.end
	let dates: string[] = []
	$: {
		dates = []
		const recordWithMostRecords = data.records.reduce((maxRecord, current) =>
			current.records.length > maxRecord.records.length ? current : maxRecord
		)
		for (let record of recordWithMostRecords.records) {
			const recordDate = record.time.substring(0, 7);
			if (recordDate < start) {
				continue
			}
			if (recordDate > end) {
				break
			}
			dates.push(recordDate)
		}
	}

	let pricesArr: (number | null)[][] = []
	$: {
		pricesArr = data.records.map(townRecord => {
			let townPrices = Array(dates.length).fill(null);
			for (let timeRecord of townRecord.records) {
				const dateIndex = dates.indexOf(timeRecord.time.substring(0, 7));
				if (dateIndex !== -1) {
					townPrices[dateIndex] = timeRecord.averageResalePrice;
				}
			}
			return townPrices;
		})
	}

	function generateColors(count: number) {
		const colors = []
		for (let i = 0; i < count; i++) {
			const hue = (i * 360) / count // even distribute hues
			const color = {
				borderColor: `hsl(${hue}, 70%, 50%)`,
				backgroundColor: `hsl(${hue}, 70%, 60%)`,
			}
			colors.push(color)
		}
		return colors
	}

	$: colors = generateColors(data.records.length)
	let generatedObjects: DataSet[] = []
	$: {
		generatedObjects = data.records.map((record, i) => ({
			label: record.town,
			data: pricesArr[i],
			xAxis: dates,
			borderColor: colors[i].borderColor,
			backgroundColor: colors[i].backgroundColor
		}))
		console.log(generatedObjects)
	}
</script>

<div class="flex flex-col justify-center items-center h-screen w-full">
	<form>
		<input type="text" class="input" name="start" bind:value={start}>
		<input type="text" class="input" name="end" bind:value={end}>
	</form>
	
	<LineChart
	{...{
		stats: generatedObjects,
		label: "Price(SGD)"
	}}>
	</LineChart>

</div>