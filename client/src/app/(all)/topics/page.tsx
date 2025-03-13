import { Topics } from "@/components/blocks/topics";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Список тем | Educ GPT",
    description: "AI-Репетитор у вас под рукой",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Список тем | Educ GPT',
        description: 'AI-Репетитор у вас под рукой',
        url: "https://educgpt.ru",
        siteName: 'Список тем | Educ GPT',
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

export default async function TopicsPage() {
    return (
        <div className="container mt-5">
            <div className="text">
                <h1>Список тем</h1>
            </div>
            <Topics />
        </div>
    );
}
