'use server'

import { baseFetch } from "../baseFetch"


export async function getDialog(id: number, token: string) {
    const data = await baseFetch("/dialogs/" + id, null, "GET", token)
    return data
}