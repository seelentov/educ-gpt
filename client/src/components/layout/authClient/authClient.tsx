'use client'

import { refresh } from "@/core/api/auth/refresh"
import { useLocalStorage } from "@/core/hooks/useLocalStorage"
import { useEffect, useRef } from "react"

export function AuthClient() {
    const [_, setToken] = useLocalStorage("token", "")
    const [refreshToken, __] = useLocalStorage("refresh_token", "")


    useEffect(() => {
        const refreshTokenAsync = async () => {
            try {
                if (refreshToken !== "") {
                    const data = await refresh(refreshToken);

                    if (data?.token) {
                        setToken(data.token)
                    }
                }
            } catch (error) {
                console.error(error)
            }
        };

        const interval = setInterval(refreshTokenAsync, 60000);

        return () => {
            if (interval) {
                clearInterval(interval);
            }
        };
    }, [refreshToken])

    return null
}