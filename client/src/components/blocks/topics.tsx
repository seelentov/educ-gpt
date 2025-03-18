'use client'

import { getTopics } from "@/core/api/roadmap/topics";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { Cards } from "../ui/cards";
import { showToast } from "../utils/toast";

export function Topics() {

    const [token] = useLocalStorage("token", "")
    const [topics, setTopics] = useState<Topic[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const router = useRouter()



    useEffect(() => {
        (async () => {
            setIsLoading(true)
            const res = await getTopics(token)
            if (res?.error) {
                console.error(res.error)
                showToast("error", res.error)
                router.refresh()
            }
            else {
                setTopics(res)
            }
            setIsLoading(false)
        })()
    }, [token, router])

    const list = useMemo(() => topics.map((t) => {
        const item = {
            slug: String(t.id),
            title: t.title,
            descInfo: ""
        }

        if (t?.scores && t.scores > 0) {
            item.descInfo = String(t.scores)
        }

        return item
    }), [topics, router])

    return <Cards linkPrefix="/topics" list={list} isLoading={isLoading} />;
}