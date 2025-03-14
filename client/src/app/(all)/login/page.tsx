import { LoginForm } from "@/components/forms/loginForm";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Вход | Educ GPT",
    description: "Вход в аккаунт Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Вход | Educ GPT',
        description: 'Вход в аккаунт Educ GPT',
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

export default function Login() {
    return (
        <div className="container">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-12 col-md-5">
                    <LoginForm />
                </div>
            </div>
        </div>
    );
}
