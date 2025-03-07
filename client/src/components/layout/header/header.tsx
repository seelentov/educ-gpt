import Link from "next/link"
import { usePathname } from "next/navigation"

export async function Header() {

    const pathname = usePathname()

    const routes = [
        {
            link: "/topics",
            title: "Список тем"
        },
        {
            link: "/leaderboard",
            title: "Таблица лидеров"
        },
        {
            link: "/donut",
            title: "Поддержкать проект"
        },

    ]
    return (
        <div className="container">
            <header className="d-flex flex-wrap align-items-center justify-content-center justify-content-md-between py-3 mb-4 border-bottom">
                <ul className="nav col-12 col-md-auto mb-2 justify-content-center mb-md-0">
                    {routes.map(({ link, title }) =>
                        <li><Link href={link} className={"nav-link px-2 " + (pathname === link ? "link-secondary" : "")}>{title}</Link></li>
                    )}

                    <li><a href="#" className="nav-link px-2 link-dark">Features</a></li>
                    <li><a href="#" className="nav-link px-2 link-dark">Pricing</a></li>
                    <li><a href="#" className="nav-link px-2 link-dark">FAQs</a></li>
                    <li><a href="#" className="nav-link px-2 link-dark">About</a></li>
                </ul>

                {(pathname != "/login" && pathname != "/signup") &&
                    <div className="col-md-3 text-end">
                        <Link href="/login" type="button" className="btn btn-outline-primary me-2">Login</Link>
                        <Link href="/signup" type="button" className="btn btn-primary">SignUp</Link>
                    </div>
                }
            </header>
        </div>
    )
}