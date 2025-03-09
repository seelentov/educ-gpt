'use client'

import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import Link from "next/link";
import { useEffect } from "react";
import { useRouter } from 'next/navigation'

export function Welcome() {

    const [token] = useLocalStorage("token", "")

    const router = useRouter()

    useEffect(() => {
        if (token !== "") {
            router.push("/topics")
        }
    })

    return (
        <div className="full-height">
            <div className="jumbotron">
                <div className="container">
                    <h1 className="display-3">AI-Репетитор у вас под рукой</h1>
                    <p>Получи знания и практику, адаптированные специально для тебя! Выбирай тему, изучай детали с помощью AI-помощника и выполняй задачи для закрепления материала. Индивидуальный подход к обучению!</p>
                    <p>
                        <Link href="/signup" className="btn btn-primary btn-lg" role="button">Присоединиться</Link>   <Link href="/login" className="btn btn-secondary btn-lg" role="button">Войти</Link>
                    </p>
                </div>
            </div>
        </div>
    )
}