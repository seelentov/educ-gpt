'use server'

import { baseFetch } from "../baseFetch"


export async function getTheme(topicId: string, themeId: string, token: string = "") {
    const data = await baseFetch("/roadmap/" + topicId + "/" + themeId, null, "GET", token)
    return data
}