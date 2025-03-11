import { Themes } from "@/components/blocks/themes";
import { getTopicInfo } from "@/core/api/roadmap/topic_info";

interface ITopicPageParams {
    params: {
        topic: number
    }
}

export default async function TopicPage({ params }: ITopicPageParams) {

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
