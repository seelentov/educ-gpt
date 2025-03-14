'use server'

import { baseFetch } from "../baseFetch"


export async function compile(code: string, language: string, token: string = "") {
    const data = await baseFetch(
        "/utils/compile",
        JSON.stringify({
            code,
            language
        }),
        "POST",
        token
    )
    return data
}