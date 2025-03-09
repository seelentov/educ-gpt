'use server'

import { baseFetch } from "../baseFetch";


export async function login(credential: string, password: string) {
    const data = await baseFetch(
        "/api/auth/login",
        JSON.stringify({
            credential,
            password,
        }),
        "POST"
    )

    return data;
}