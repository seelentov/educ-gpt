'use server'

import { baseFetch } from "../baseFetch";


export async function changeEmailActivate(key: string, userId: string) {
    const data = await baseFetch(
        "/auth/change_email/" + key + "/" + userId,
        null,
        "POST"
    )

    return data;
}