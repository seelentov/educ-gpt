'use server'

import { baseFetch } from "../baseFetch"


export async function getProblems(topicId: number, themeId: number, token: string = "") {
    const data = await baseFetch("/roadmap/problems" + topicId + "/" + themeId, null, "GET", token)
    return data
}