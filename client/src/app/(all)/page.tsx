import { Welcome } from "@/components/blocks/welcome";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Educ GPT",
    description: "AI-Репетитор у вас под рукой",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Educ GPT',
        description: 'AI-Репетитор у вас под рукой',
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

export default function Home() {
    return (
        <Welcome />
    );
}
