'use server'

import { baseFetch } from "../baseFetch";


export async function activate(key: string) {
    const data = await baseFetch(
        "/auth/activate/" + key,
        null,
        "POST"
    )

    return data;
}