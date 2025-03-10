'use client'

import { useRouter } from "next/navigation"
import { useEffect } from "react"

interface IRedirectorProps {
    timeout: number
    link: string
}

export function Redirector({ timeout, link }: IRedirectorProps) {
    const router = useRouter()

    useEffect(() => {
        const timeOut = setTimeout(() => {
            router.replace(link)
        }, timeout)
    }, [])

    return null
}