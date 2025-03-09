import { Topics } from "@/components/blocks/topics";

export default async function TopicsPage() {
    return (
        <div className="container mt-5">
            <div className="text">
                <h1>Список тем</h1>
            </div>
            <Topics />
        </div>
    );
}
