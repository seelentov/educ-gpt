'use client'

import { refresh } from "@/core/api/auth/refresh"
import { useLocalStorage } from "@/core/hooks/useLocalStorage"
import { useEffect } from "react"

export function AuthClient() {
    const [, setToken] = useLocalStorage("token", "")


    const getRefreshToken = () => {
        const value = window.localStorage.getItem("refresh_token")

        const v = value ? JSON.parse(value) : ""

        return String(v)
    }


    useEffect(() => {
        const refreshTokenAsync = async () => {
            try {

                const token = getRefreshToken()

                if (token !== "") {
                    const data = await refresh(token);

                    if (data?.token) {
                        setToken(data.token)
                    }
                    else {
                        await refreshTokenAsync()
                    }
                }
            } catch (error) {
                console.error(error)
            }
        };

        refreshTokenAsync()

        const interval = setInterval(refreshTokenAsync, 60000);

        return () => {
            if (interval) {
                clearInterval(interval);
            }
        };
    }, [setToken])

    return null
}