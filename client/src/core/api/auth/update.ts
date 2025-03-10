'use server'

import { baseFetch } from "../baseFetch";


export async function updateUser(name: string, number: string, chat_gpt_token: string, chat_gpt_model: string, avatar_url: string, token: string) {
    const data = await baseFetch(
        "/auth/update",
        JSON.stringify({
            name,
            chat_gpt_model,
            number,
            avatar_url,
            chat_gpt_token
        }),
        "POST",
        token
    )

    return data;
}