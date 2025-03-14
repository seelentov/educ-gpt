'use server'

import { baseFetch } from "../baseFetch"


export async function resolve(problem_id: number, answer: string, language: string, token: string = "") {
    const data = await baseFetch(
        "/roadmap/resolve",
        JSON.stringify({
            problem_id,
            answer,
            language
        }),
        "POST",
        token
    )
    return data
}