import { createClient } from "@libsql/client"
import { drizzle } from "drizzle-orm/libsql"
import { DB_TOKEN, DB_URL } from "$env/static/private"

const client = createClient({ url: DB_URL, authToken: DB_TOKEN })
export const db = drizzle(client, { schema: undefined })