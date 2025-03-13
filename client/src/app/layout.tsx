import 'bootstrap/dist/css/bootstrap.css';
import '@/styles/styles.scss'

import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import React from "react";
import { AuthClient } from '@/components/layout/authClient/authClient';


const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

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

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {

    return (
        <html lang="ru">
            <body className={`${geistSans.variable} ${geistMono.variable}`}>
                <AuthClient />
                <main>
                    {children}
                </main>
            </body>
        </html>
    );
}
