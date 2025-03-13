import { ChangeEmailForm } from "@/components/forms/changeEmailForm";
import { ChangePasswordForm } from "@/components/forms/changePasswordForm";
import EditProfileForm from "@/components/forms/editProfileForm";
import { Metadata } from "next";


export const metadata: Metadata = {
    title: "Профиль | Educ GPT",
    metadataBase: new URL("https://educgpt.ru"),
    openGraph: {
        title: 'Профиль | Educ GPT',
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

export default function Profile() {
    return (
        <div className="container mt-5">
            <div className="text">
                <h1>Настройки профиля</h1>
            </div>
            <EditProfileForm />
            <hr />
            <div className="text">
                <h2>Сменить пароль</h2>
            </div>
            <ChangePasswordForm />
            <hr />
            <div className="text">
                <h2>Сменить почту</h2>
            </div>
            <ChangeEmailForm />
        </div>
    );
}
