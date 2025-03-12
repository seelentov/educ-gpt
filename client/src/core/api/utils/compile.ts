'use server'

import { baseFetch } from "../baseFetch"


export async function compile(code: string, token: string = "") {
    const data = await baseFetch(
        "/utils/compile",
        JSON.stringify({
            code
        }),
        "POST",
        token
    )
    return data
}