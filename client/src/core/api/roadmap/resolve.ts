'use server'

import { baseFetch } from "../baseFetch"


export async function resolve(problem_id: number, answer: string, token: string = "") {
    const data = await baseFetch(
        "/roadmap/resolve",
        JSON.stringify({
            problem_id,
            answer
        }),
        "POST",
        token
    )
    return data
}