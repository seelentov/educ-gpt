'use server'

import { baseFetch } from "../baseFetch"


export async function checkAnswerUtil(problem: string, answer: string, language: string, token: string = "") {
    const data = await baseFetch(
        "/utils/check_answer",
        JSON.stringify({
            problem,
            language,
            answer
        }),
        "POST",
        token
    )
    return data
}