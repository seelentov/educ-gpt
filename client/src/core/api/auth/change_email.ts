'use server'

import { baseFetch } from "../baseFetch";


export async function changeEmail(email: string, token: string) {
    const data = await baseFetch(
        "/auth/change_email/task",
        JSON.stringify({
            email
        }),
        "POST",
        token
    )

    return data;
}