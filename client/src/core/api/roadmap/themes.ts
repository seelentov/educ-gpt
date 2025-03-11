'use server'

import { baseFetch } from "../baseFetch"


export async function getThemes(id: number, token: string = "") {
    const data = await baseFetch("/roadmap/" + id, null, "GET", token)
    return data
}