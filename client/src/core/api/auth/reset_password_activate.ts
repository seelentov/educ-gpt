'use server'

import { baseFetch } from "../baseFetch";


export async function resetPasswordActivate(key: string, userId: string, password: string) {
    const data = await baseFetch(
        "/auth/reset/" + key + "/" + userId,
        JSON.stringify({
            password
        }),
        "POST"
    )

    return data;
}