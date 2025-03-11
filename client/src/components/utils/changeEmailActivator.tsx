'use client'

import { changeEmailActivate } from "@/core/api/auth/change_email_activate"
import { useLocalStorage } from "@/core/hooks/useLocalStorage"
import { Redirector } from "./redirector"
import { useEffect, useState } from "react"

interface IChangeEmailActivateProps {
    key: string
}

export default function ChangeEmailActivator({ key }: IChangeEmailActivateProps) {
    const [token] = useLocalStorage("token", "")

    const [data, setData] = useState({
        status: "",
        error: "",
    })

    useEffect(() => {
        (async () => {
            console.log(key)
            console.log(token)
            if (token !== "") {
                const res = await changeEmailActivate(key, token)
                setData(res)
            }
        })()

    }, [token, key])



    return <div className="container mt-5" >
        <p>{data?.status || data?.error
        } </p>
        < hr />
    </div>
}