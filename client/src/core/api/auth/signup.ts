'use server'

import { baseFetch } from "../baseFetch";


export async function signup(name: string, email: string, password: string) {
    const data = await baseFetch(
        "/auth/register",
        JSON.stringify({
            name,
            email,
            password
        }),
        "POST"
    )

    return data;
}