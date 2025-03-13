import { SignUpForm } from "@/components/forms/signUpForm";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Регистрация | Educ GPT",
    description: "Регистрация аккаунта Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Регистрация | Educ GPT',
        description: 'Регистрация аккаунта Educ GPT',
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

export default function Signup() {
    return (
        <div className="container">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-5">
                    <SignUpForm />
                </div>
            </div>
        </div>
    );
}
