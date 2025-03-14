'use server'

import { baseFetch } from "../baseFetch";


export async function signup(name: string, email: string, password: string, chat_gpt_token: string) {
    const data = await baseFetch(
        "/auth/register",
        JSON.stringify({
            name,
            email,
            password,
            chat_gpt_token
        }),
        "POST"
    )

    return data;
}