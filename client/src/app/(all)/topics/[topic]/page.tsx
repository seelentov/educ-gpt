import { Themes } from "@/components/blocks/themes";
import { getTopicInfo } from "@/core/api/roadmap/topic_info";
import { Metadata } from "next";
import { notFound } from "next/navigation";

export async function generateMetadata(
    props: { params: Params }
): Promise<Metadata> {

    const params = await props.params;
    const topic: Topic = await getTopicInfo(params.topic)

    if (!topic) {
        notFound();
    }

    return {
        title: `${topic.title} | Educ GPT`,
        description: `Разделы по теме ${topic.title}`,
        openGraph: {
            title: `${topic.title} | Educ GPT`,
            description: `Разделы по теме ${topic.title}`,
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

type Params = Promise<{ topic: number }>

export default async function TopicPage(props: { params: Params }) {
    const params = await props.params;
    const topic = await getTopicInfo(params.topic)

    return (
        <div className="container mt-5">
            <div className="text">
                <h1>{topic.title}</h1>
            </div>
            <Themes id={topic.id} />
        </div>
    );
}
