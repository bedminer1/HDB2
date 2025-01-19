<script lang="ts">
	import LineChart from "$lib/ components/LineChart.svelte";
	export let data: {
		records: townRecords[]
	}
	
	const dates: string[] = []
	for (let record of data.records[0].records) {
		dates.push(record.time.substring(0, 10))
	}
	const pricesArr: number[][] = []
	for (let townRecord of data.records) {
		let townPrices: number[] = []
		for (let timeRecord of townRecord.records) {
			townPrices.push(timeRecord.averageResalePrice)
		}
		pricesArr.push(townPrices)
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

	const colors = generateColors(data.records.length)
	const generatedObjects: DataSet[] = []
	for (let i = 0; i < data.records.length; i++) {
		const record = data.records[i]
		const obj = {
			label: record.town,
			data: pricesArr[i],
			xAxis: dates,
			borderColor: colors[i].borderColor,
			backgroundColor: colors[i].backgroundColor
		}
		generatedObjects.push(obj)
	}

</script>

<div class="flex flex-col justify-center items-center h-screen w-full">
	<LineChart
	{...{
		stats: generatedObjects,
		label: "Price(SGD)"
	}}>

	</LineChart>
</div>