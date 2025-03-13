import { notFound } from "next/navigation";
import { Metadata } from "next";
import { getThemeInfo } from "@/core/api/roadmap/theme_info";
import { Editor } from "@/components/blocks/editor";

type Params = Promise<{ topic: number, theme: number }>

export async function generateMetadata(
    props: { params: Params }
): Promise<Metadata> {

    const params = await props.params;
    const theme: Theme = await getThemeInfo(params.theme)

    if (!theme) {
        notFound();
    }

    return {
        title: `${theme.topic?.title}. ${theme.title} | Educ GPT`,
        description: `Теория и практика по теме ${theme.topic?.title}. ${theme.title}.`,
        openGraph: {
            title: `${theme.topic?.title}. ${theme.title} | Educ GPT`,
            description: `Теория и практика по теме ${theme.topic?.title}. ${theme.title}.`,
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
    }
}

export default function ThemePage() {
    <Editor />
}
