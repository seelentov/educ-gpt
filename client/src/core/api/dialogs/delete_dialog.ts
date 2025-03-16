'use server'

import { baseFetch } from "../baseFetch"


export async function deleteDialog(id: number, token: string) {
    const data = await baseFetch("/dialogs/" + id, null, "DELETE", token)
    return data
}