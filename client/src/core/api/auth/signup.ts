'use server'

import { baseFetch } from "../baseFetch";


export async function signup(name: string, email: string, number: string, password: string, chat_gpt_token: string) {
    const data = await baseFetch(
        "/auth/register",
        JSON.stringify({
            name,
            email,
            number,
            password,
            chat_gpt_token
        }),
        "POST"
    )

    return data;
}