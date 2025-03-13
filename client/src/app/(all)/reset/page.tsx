import { ResetPasswordTaskForm } from "@/components/forms/resetPasswordTaskForm";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Восстановление аккаунта | Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Восстановление аккаунта | Educ GPT',
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

export default function ResetPassword() {
    return (
        <div className="conteiner mt-5">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-12 col-md-5">
                    <ResetPasswordTaskForm />
                </div>
            </div>
        </div>
    )
}