'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect, useMemo, useState } from 'react'
import Image from "next/image";
import { useLocalStorage } from '@/core/hooks/useLocalStorage'
import { me } from '@/core/api/auth/me'
import { HOST_URL_PROD } from '@/core/api/api';

export function Header() {

    const pathname = usePathname()

    const router = useRouter()

    const [token, setToken] = useLocalStorage("token", "")
    const [, setRefreshToken] = useLocalStorage("refresh_token", "")

    const [isLogged, setIsLogged] = useState<boolean>(false)

    const [avatarUrl, setAvatarUrl] = useState<string>("")

    useEffect(() => {
        (async () => {

            const data = await me(token)

            if (data?.avatar_url) {
                setAvatarUrl(HOST_URL_PROD + data.avatar_url)
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
            console.error(error)
        }

    }, [pathname, token])

    const logout = (e: any) => {
        e.preventDefault()

        setToken("")
        setRefreshToken("")
        router.push("/")
    }

    const routes = useMemo(() => [
        {
            link: "/topics",
            title: "Список тем",
        },

    ], [])

    return (
        <div className='header'>
            <div className="container">
                <header className="d-flex align-items-center justify-content-center justify-content-md-between py-sm-3 py-1 border-bottom">
                    <ul className="nav col-sm-12 col-md-auto col-7 mb-sm-2 justify-content-start justify-content-sm-center mb-md-0">
                        {routes.map(({ link, title }) =>
                            <li key={link}><Link href={link} className={"nav-link px-2 " + (pathname === link ? "link-secondary disabled" : "link-dark")}>
                                {title}
                            </Link></li>
                        )}
                        <li>
                            <a target='_blank' href="/swagger/index.html" className="nav-link px-2 link-dark">
                                API
                            </a>
                        </li>
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
                                        ? <Image style={{ objectFit: 'cover' }} src={avatarUrl} alt="" width={32} height={32} className="rounded-circle" />
                                        : <Image style={{ objectFit: 'cover' }} src={"/misc/empty_avatar.jpg"} alt="" width={32} height={32} className="rounded-circle" />
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
