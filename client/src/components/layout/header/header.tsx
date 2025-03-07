'use client'

import { useWidth } from '@/core/hooks/useWidth'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useMemo } from 'react'
import Image from "next/image";

export function Header() {

    const pathname = usePathname()

    const routes = useMemo(() => [
        {
            link: "/topics",
            title: "Список тем",
            icon: "/book.svg"
        },
        {
            link: "/leaderboard",
            title: "Таблица лидеров",
            icon: "/person.svg"
        },
        {
            link: "/donut",
            title: "Поддержкать проект",
            icon: "/donat.svg"
        },

    ], [])

    const width = useWidth()

    return (
        <div className='header'>
            <div className="container">
                <header className="d-flex flex-wrap align-items-center justify-content-center justify-content-md-between py-3 mb-4 border-bottom">
                    <ul className="nav col-12 col-md-auto mb-2 justify-content-center mb-md-0">
                        {routes.map(({ link, title, icon }) =>
                            <li><Link href={link} className={"nav-link px-2 " + (pathname === link ? "link-secondary disabled" : "link-dark")}>
                                {
                                    width > 720 ? title : <Image src={icon} alt={""} width={20} height={20} />
                                }
                            </Link></li>
                        )}
                    </ul>
                    {(pathname != "/login" && pathname != "/signup") &&
                        <div className="col-md-3 text-end">
                            <Link href="/login" type="button" className="btn btn-outline-primary me-2">Login</Link>
                            <Link href="/signup" type="button" className="btn btn-primary">SignUp</Link>
                        </div>
                    }
                </header>
            </div>
        </div>
    )
}
