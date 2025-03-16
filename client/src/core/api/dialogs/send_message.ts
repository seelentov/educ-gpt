'use server'

import { baseFetch } from "../baseFetch"


export async function sendMessage(id: number, message: string, token: string) {
    const data = await baseFetch(
        "/dialogs/" + id,
        JSON.stringify({
            message
        }),
        "POST",
        token
    )
    return data
}