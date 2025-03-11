'use server'

import { baseFetch } from "../baseFetch"


export async function getTopicInfo(id: number) {
    const data = await baseFetch("/roadmap/" + id + "/info", null, "GET")
    return data
}