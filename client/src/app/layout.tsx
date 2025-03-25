'use client'

import 'bootstrap/dist/css/bootstrap.css';
import '@/styles/styles.scss'

import { Geist, Geist_Mono } from "next/font/google";
import React from "react";
import { Footer } from "@/components/layout/footer/footer";
import { AuthClient } from '@/components/layout/authClient/authClient';
import { Chat } from '@/components/layout/chat/chat';
import ToastProvider from '@/components/providers/toast';
import { Header } from '@/components/layout/header/header';
import { usePathname } from 'next/navigation';


const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const pathname = usePathname()

    const isEditor = pathname.match(/^\/topics\/\d+\/\d+$/) != null

    return (
        <html lang="ru">
            <body className={`${geistSans.variable} ${geistMono.variable}`}>
                <ToastProvider>
                    <AuthClient />
                    {!isEditor && <Header />}
                    {isEditor
                        ? <main>
                            {children}
                        </main>
                        : <div className='wrapper'>
                            <main>
                                {children}
                            </main>
                        </div>
                    }
                    {!isEditor && <Footer />}
                    <Chat />
                </ToastProvider>
            </body>
        </html >
    );
}
