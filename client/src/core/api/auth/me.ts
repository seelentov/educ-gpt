'use server'

import { baseFetch } from "../baseFetch"


export async function me(token: string) {
    const data = await baseFetch("/auth/me", null, "POST", token)
    return data
}