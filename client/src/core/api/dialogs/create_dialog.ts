'use server'

import { baseFetch } from "../baseFetch"


export async function createDialog(token: string) {
    const data = await baseFetch("/dialogs/", null, "POST", token)
    return data
}