import Image from "next/image";

export function Footer() {
    return (
        <div className="container">
            <footer className="d-flex flex-wrap justify-content-between align-items-center py-3 my-4 border-top">
                <ul className="nav col-12 justify-content-end list-unstyled d-flex">
                    <li className="ms-3">
                        <a className="text-muted" target="_blank" href="https://github.com/seelentov/">
                            <Image src="/icons/github.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                    <li className="ms-3">
                        <a className="text-muted" target="_blank" href="https://t.me/komkov01">
                            <Image src="/icons/telegram.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                    <li className="ms-3">
                        <a className="text-muted" target="_blank" href="https://hh.ru/resume/ffae0a05ff0b8281f20039ed1f587964334e32">
                            <Image src="/icons/hh.svg" alt={""} width={40} height={40} />
                        </a>
                    </li>
                </ul>
            </footer>
        </div>

    )
}