import path from 'path'
import { readFile } from 'fs/promises'

export async function load() {
    const filePath = path.resolve("./src/lib/town_query_data.json")

    try {
        const fileData = await readFile(filePath, 'utf-8')
        const queryData = JSON.parse(fileData)
        let towns: string[] = queryData.towns
        let start: string = queryData.start
        let end: string = queryData.end

        let townsQueries = ""
        for (let town of towns) {
            townsQueries += "&towns=" + town
        }
        
        const apiURL = `http://localhost:4000/town_stats?start=${start}&end=${end}&flatType=3%20ROOM` + townsQueries
        const response = await fetch(apiURL)
        if (!response.ok) {
            throw new Error(`API request failed with status ${response.status}`) 
        }

        const data = await response.json()
        let { records }: { records: timeBasedRecord[]} = data
        return {
            records
        }

    } catch (err) {
        console.log(err)
    }
}