'use client'

import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { Cards } from "../ui/cards";
import { getThemes } from "@/core/api/roadmap/themes";
import { showToast } from "../utils/toast";

interface IThemesProps {
    id: number
}

export function Themes({ id }: IThemesProps) {

    const [token] = useLocalStorage("token", "")
    const [themes, setThemes] = useState<Theme[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const router = useRouter()

    useEffect(() => {
        if (token === "") {
            router.replace("/login")
        }
    })

    useEffect(() => {
        (async () => {
            setIsLoading(true)
            const res = await getThemes(id, token)
            if (res?.error) {
                console.error(res.error)
                showToast("error", res.error)
                router.refresh()
            }
            else {
                setThemes(res)
            }
            setIsLoading(false)
        })()
    }, [token, id, router])

    const list = useMemo(() => themes.map((t) => {
        const item = {
            slug: t.id,
            title: t.title,
            descInfo: ""
        }

        if (t?.scores && t.scores > 0) {
            item.descInfo = String(t.scores)
        }

        return item
    }), [themes])

    return <Cards linkPrefix={"/topics/" + id} list={list} isLoading={isLoading} loadingText="Подбираем список тем, учитывая вашу статистику" />;
}