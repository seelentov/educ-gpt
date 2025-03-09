'use server'

import { baseFetch } from "../baseFetch";


export async function refresh(refresh_token: string) {
    const data = await baseFetch(
        "/api/auth/refresh",
        JSON.stringify({
            refresh_token,
        }),
        "POST"
    )

    return data;
}