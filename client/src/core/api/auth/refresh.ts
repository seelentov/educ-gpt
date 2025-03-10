'use server'

import { baseFetch } from "../baseFetch";


export async function refresh(refresh_token: string) {
    const data = await baseFetch(
        "/auth/refresh",
        JSON.stringify({
            refresh_token,
        }),
        "POST"
    )

    return data;
}