'use server'

import { baseFetch } from "../baseFetch";


export async function resetPassword(credential: string) {
    const data = await baseFetch(
        "/auth/reset/task",
        JSON.stringify({
            credential,
        }),
        "POST"
    )

    return data;
}