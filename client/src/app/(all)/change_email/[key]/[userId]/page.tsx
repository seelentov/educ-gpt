import { changeEmailActivate } from "@/core/api/auth/change_email_activate";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Изменить E-mail | Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Изменить E-mail | Educ GPT',
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

type Params = Promise<{ key: string, userId: string }>

export default async function ChangeEmailActivatePage(props: { params: Params }) {
    const params = await props.params;
    const data = await changeEmailActivate(params.key, params.userId)

    return <div className="container mt-5">
        <p>{data?.status || data?.error}</p>
        <hr />
    </div>
}