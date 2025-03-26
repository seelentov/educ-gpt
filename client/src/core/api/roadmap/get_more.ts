'use server'

import { baseFetch } from "../baseFetch"


export async function getMoreInfo(topic: string, theme: string, messages: string[]) {
    const data = await baseFetch(
        "/roadmap/more/" + topic + "/" + theme,
        JSON.stringify({
            messages
        }),
        "POST"
    )
    return data
}