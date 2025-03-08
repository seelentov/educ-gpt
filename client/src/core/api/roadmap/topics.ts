'use server'

import { baseFetch } from "../baseFetch"


export async function getTopics(token: string = "") {
    const data = await baseFetch("/api/roadmap", null, "GET", token)
    return data
}