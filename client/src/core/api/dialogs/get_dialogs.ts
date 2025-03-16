'use server'

import { baseFetch } from "../baseFetch"


export async function getDialogs(token: string) {
    const data = await baseFetch("/dialogs/", null, "GET", token)
    return data
}