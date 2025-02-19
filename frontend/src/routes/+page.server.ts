import path from 'path'
import { readFile } from 'fs/promises'
import { superValidate } from 'sveltekit-superforms'
import { formSchema } from './schema'
import { zod } from "sveltekit-superforms/adapters"

export async function load() {
    const filePath = path.resolve("./src/lib/town_query_data.json")

    const fileData = await readFile(filePath, 'utf-8')
    const queryData = JSON.parse(fileData)
    let towns: string[] = queryData.towns
    let start: string = queryData.start
    let end: string = queryData.end

    let townsQueries = ""
    for (let town of towns) {
        townsQueries += "&towns=" + town.replace(/ /g, "%20")
    }
    
    const apiURL = `http://localhost:4000/town_stats?start=${start}&end=${end}&flatType=3%20ROOM` + townsQueries
    const response = await fetch(apiURL)
    if (!response.ok) {
        console.error(`API request failed with status ${response.status}`) 
    }

    const data = await response.json()

    const form = await superValidate(zod(formSchema))
    form.data.start = start
    form.data.end = end

    let { records }: { records: townRecords[]} = data
    return {
        records,
        start: start + "-01",
        end: end + "-01",
        form
    }
}