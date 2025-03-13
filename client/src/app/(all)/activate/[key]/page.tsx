import { Redirector } from "@/components/utils/redirector"
import { activate } from "@/core/api/auth/activate"
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Активация | Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Активация | Educ GPT',
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

type Params = Promise<{ key: string }>

export default async function Activate(props: { params: Params }) {
    const params = await props.params;
    const data = await activate(params.key)

    return <div className="container mt-5">
        <p>{data?.status || data?.error}</p>
        <hr />
        {!data?.error && <Redirector timeout={5000} link="/login" />}
    </div>
}