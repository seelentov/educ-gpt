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
            icon: "/icons/book.svg"
        },
        {
            link: "/leaderboard",
            title: "Таблица лидеров",
            icon: "/icons/person.svg"
        },
        {
            link: "/donut",
            title: "Поддержать проект",
            icon: "/icons/donat.svg"
        },

    ], [])

    const width = useWidth()

    return (
        <div className='header'>
            <div className="container">
                <header className="d-flex align-items-center justify-content-center justify-content-md-between py-sm-3 py-1 border-bottom">
                    <ul className="nav col-sm-12 col-md-auto col-7 mb-sm-2 justify-content-start justify-content-sm-center mb-md-0">
                        {routes.map(({ link, title, icon }) =>
                            <li key={link}><Link href={link} className={"nav-link px-2 " + (pathname === link ? "link-secondary disabled" : "link-dark")}>
                                {
                                    width > 720 ? title : <Image src={icon} alt={""} width={30} height={30} />
                                }
                            </Link></li>
                        )}
                    </ul>
                    {(pathname != "/login" && pathname != "/signup") &&
                        <div className="col-md-3 text-end col-5">
                            <Link href="/login" type="button" className="btn btn-outline-primary btn-sm me-2">Login</Link>
                            <Link href="/signup" type="button" className="btn btn-primary btn-sm">SignUp</Link>
                        </div>
                    }
                </header>
            </div>
        </div>
    )
}
