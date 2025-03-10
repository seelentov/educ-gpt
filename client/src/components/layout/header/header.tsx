'use client'

import { useWidth } from '@/core/hooks/useWidth'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect, useMemo, useState } from 'react'
import Image from "next/image";
import { useLocalStorage } from '@/core/hooks/useLocalStorage'
import { me } from '@/core/api/auth/me'
import { HOST_URL } from '@/core/api/api'

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
    const router = useRouter()

    const [token, setToken] = useLocalStorage("token", "")
    const [__, setRefreshToken] = useLocalStorage("refresh_token", "")

    const [isLogged, setIsLogged] = useState<boolean>(false)

    const [avatarUrl, setAvatarUrl] = useState<string>("")

    useEffect(() => {
        (async () => {

            const data = await me(token)

            if (data?.avatar_url) {
                setAvatarUrl(data.avatar_url)
            }

        })()
    }, [pathname, token, isLogged])

    useEffect(() => {
        try {
            const value = window.localStorage.getItem("token")

            const v = value ? JSON.parse(value) : ""

            if (v !== "") {
                setIsLogged(true)
            } else {
                setIsLogged(false)
            }

        } catch (error) {
            console.log(error)
        }

    }, [pathname, token])

    const logout = (e: any) => {
        e.preventDefault()

        setToken("")
        setRefreshToken("")
        router.push("/")
    }

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
                    {(pathname != "/login" && pathname != "/signup" && !isLogged) &&
                        <div className="col-md-3 text-end col-5">
                            <Link href="/login" type="button" className="btn btn-outline-primary btn-sm me-2 d-none d-sm-inline-block">Войти</Link>
                            <Link href="/signup" type="button" className="btn btn-primary btn-sm">Присоединиться</Link>
                        </div>
                    }
                    {isLogged &&
                        <div className="col-md-3 text-end col-5 d-flex align-items-start justify-content-end">
                            <Link href="/profile" type="button" className="me-2">
                                {
                                    avatarUrl && avatarUrl !== ""
                                        ? <Image src={HOST_URL + "/storage/" + avatarUrl} alt="" width={32} height={32} className="rounded-circle" />
                                        : <Image src={"/misc/empty_avatar.jpg"} alt="" width={32} height={32} className="rounded-circle" />
                                }

                            </Link>
                            <button onClick={logout} type="button" className="btn btn-outline-primary btn-sm">Выйти</button>
                        </div>
                    }
                </header>
            </div>
        </div>
    )
}
