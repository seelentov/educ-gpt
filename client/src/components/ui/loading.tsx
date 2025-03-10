import { BeatLoader } from "react-spinners";

interface ILoadingProps {
    min?: boolean
    color?: string
}

export function Loading({ min, color = "#0d6efd" }: ILoadingProps) {
    return <BeatLoader size={min ? 15 : 60} color={color} />
}