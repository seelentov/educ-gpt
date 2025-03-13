import { Metadata } from "next";

export const metadata: Metadata = {
    title: "404 | Educ GPT",
    description: "Страница не найдена",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: '404 | Educ GPT',
        description: 'Страница не найдена',
        url: "https://educgpt.ru",
        siteName: 'Educ GPT',
        images: [
            {
                url: '/favicon.png',
                width: 500,
                height: 500,
            },
        ],
        locale: 'ru',
        type: 'website',
    },
    icons: {
        icon: '/favicon.png',
        shortcut: '/favicon-16x16.png',
        apple: '/apple-touch-icon.png',
    },
};

export default async function NotFoundPage() {

    return (
        <div className="container mt-5">
            <div className="text col-12 d-flex justify-content-center align-items-center gap-5">
                <h1>404</h1>
                <p>Страница не найдена</p>
            </div>
        </div>
    );
}