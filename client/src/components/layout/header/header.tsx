'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useMemo } from 'react'

interface IRoute {
    title: string
    link: string
}

export function Header() {

    const routes: IRoute[] = useMemo(() => ([
        {
            title: "Темы",
            link: "/topics"
        },
        {
            title: "Профиль",
            link: "/profile"
        },
        {
            title: "Поддержка",
            link: "/donut"
        },
    ]), [])

    const pathname = usePathname()

    return (
        <header className="d-flex justify-content-center py-3">
            <ul className="nav nav-pills">
                {routes.map(({ title, link }) =>
                    <li className="nav-item">
                        <Link href="#" className={"nav-link " + pathname.includes(link) ? "active" : ""}>{title}</Link>
                    </li>
                )}
            </ul>
        </header>
        <header className="d-flex flex-wrap align-items-center justify-content-center justify-content-md-between py-3 mb-4 border-bottom">
        <a href="/" className="d-flex align-items-center col-md-3 mb-2 mb-md-0 text-dark text-decoration-none">
          <svg className="bi me-2" width="40" height="32" role="img" aria-label="Bootstrap"><use xlink:href="#bootstrap"></use></svg>
        </a>
  
        <ul className="nav col-12 col-md-auto mb-2 justify-content-center mb-md-0">
          <li><a href="#" className="nav-link px-2 link-secondary">Home</a></li>
          <li><a href="#" className="nav-link px-2 link-dark">Features</a></li>
          <li><a href="#" className="nav-link px-2 link-dark">Pricing</a></li>
          <li><a href="#" className="nav-link px-2 link-dark">FAQs</a></li>
          <li><a href="#" className="nav-link px-2 link-dark">About</a></li>
        </ul>
  
        <div className="col-md-3 text-end">
          <button type="button" className="btn btn-outline-primary me-2">Login</button>
          <button type="button" className="btn btn-primary">Sign-up</button>
        </div>
      </header>

    )
}