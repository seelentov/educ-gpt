'use server'

import { baseFetch } from "../baseFetch";


export async function updateUser(form: FormData, token: string) {
    const data = await baseFetch(
        "/auth/update",
        form,
        "PATCH",
        token,
        null,
        true
    )

    return data;
}