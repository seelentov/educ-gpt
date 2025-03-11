'use client'

import { Loading } from "@/components/ui/loading";
import { getTopics } from "@/core/api/roadmap/topics";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { Cards } from "../ui/cards";

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
                router.refresh()
            }
            else {
                setTopics(res)
            }
            setIsLoading(false)
        })()
    }, [token])

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
    }), [topics])

    return <Cards linkPrefix="/topics" list={list} isLoading={isLoading} />;
}