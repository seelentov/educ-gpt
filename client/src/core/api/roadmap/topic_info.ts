'use server'

import { baseFetch } from "../baseFetch"


export async function getTopicInfo(id: number) {
    const data = await baseFetch("/roadmap/info/topic/" + id, null, "GET")
    return data
}