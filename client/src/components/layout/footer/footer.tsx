import Link from "next/link";
import Image from "next/image";

export function Footer() {
    return (
        <div className="container">
            <footer className="d-flex flex-wrap justify-content-between align-items-center py-3 my-4 border-top">
                <div className="col-md-4 d-flex align-items-center">
                    <Link href="/" className="mb-3 me-2 mb-md-0 text-muted text-decoration-none lh-1">
                        <Image src="/icons/logo.svg" alt={""} width={30} height={30} />
                    </Link>
                    <span className="text-muted">© 2021 Company, Inc</span>
                </div>

                <ul className="nav col-md-4 justify-content-end list-unstyled d-flex">
                    <li className="ms-3">
                        <a className="text-muted" href="https://github.com/seelentov/">
                            <Image src="/icons/github.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                    <li className="ms-3">
                        <a className="text-muted" href="https://t.me/komkov01">
                            <Image src="/icons/telegram.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                    <li className="ms-3">
                        <a className="text-muted" href="https://hh.ru/resume/ffae0a05ff0b8281f20039ed1f587964334e32">
                            <Image src="/icons/hh.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                </ul>
            </footer>
        </div>

    )
}