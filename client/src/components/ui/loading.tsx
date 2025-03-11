import { BeatLoader } from "react-spinners";

interface ILoadingProps {
    min?: boolean
    color?: string
    text?: string
}

export function Loading({ min, color = "#0d6efd", text = "" }: ILoadingProps) {
    return <div className="d-flex flex-column gap-5 align-items-center">
        <BeatLoader size={min ? 15 : 60} color={color} />
        {text && <p>{text}</p>}
    </div>

}