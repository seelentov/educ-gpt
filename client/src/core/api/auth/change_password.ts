



'use server'

import { baseFetch } from "../baseFetch";


export async function changePassword(old_password: string, password: string, token: string) {
    const data = await baseFetch(
        "/auth/change_password",
        JSON.stringify({
            old_password,
            password,
        }),
        "POST",
        token
    )

    return data;
}