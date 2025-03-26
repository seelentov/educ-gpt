'use server'

import { baseFetch } from "../baseFetch"


export async function getThemeInfo(id: number) {
    const data = await baseFetch("/roadmap/info/theme/" + id, null, "GET")
    return data
}