'use client'

import { Loading } from "@/components/ui/loading";
import { getTopics } from "@/core/api/roadmap/topics";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import Link from "next/link";
import { useEffect, useState } from "react";

export function Topics() {

    const [token] = useLocalStorage("token", "")
    const [topics, setTopics] = useState<Topic[]>([])

    useEffect(() => {
        (async () => {
            const res = await getTopics(token)
            if (res?.error) {
                alert(res.error)
            }
            else {
                setTopics(res)
            }

        })()
    }, [token])

    return (
        <div className="card-desk d-flex flex-wrap">
            {(topics && topics.length > 0) && topics.map(t =>
                <Link key={t.id} className="card col-6 col-md-3" href={`/topics/${t.id}`}>
                    <div className="card-body">
                        <h2 className="card-title">{t.title}</h2>
                        {t?.scores && <small className="text-muted">Очки: {t.scores}</small>}
                    </div>
                </Link>
            )}
            {(!topics || topics.length === 0) &&
                <div className="col-12 d-flex justify-content-center align-items-center" style={{ height: 200 }}>
                    <Loading />
                </div>
            }

        </div>
    );
}