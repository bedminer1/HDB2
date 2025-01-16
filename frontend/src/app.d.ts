// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
	// interface Locals {}
	// interface PageData {}
	// interface Error {}
	// interface Platform {}
}

interface townRecords {
	town: string
	records: timeBasedRecord[]
}

interface timeBasedRecord {
	averagePricePerArea: number
	averageResalePrice: number
	numberOfUnits: number
	time: string
}

interface DataSet {
	label: string;
	data: number[];
	xAxis: string[];
	borderColor: string;
	backgroundColor: string;
  }
  