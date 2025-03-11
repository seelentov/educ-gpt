import Link from "next/link"
import { Loading } from "./loading"

interface ICardsProps {
    list: IListItem[]
    linkPrefix?: string
    isLoading?: boolean
    loadingText?: string
}

interface IListItem {

    slug: string | number
    title: string
    descInfo?: string
}


export function Cards({ list, linkPrefix = "", isLoading = false, loadingText = "" }: ICardsProps) {
    return (
        <div className="card-desk d-flex flex-wrap">
            {!isLoading && list.map(t =>
                <Link key={t.slug} className="card col-6 col-md-3" href={`${linkPrefix}/${t.slug}`}>
                    <div className="card-body">
                        <h2 className="card-title">{t.title}</h2>
                        {t?.descInfo && <small className="text-muted">Очки: {t.descInfo}</small>}
                    </div>
                </Link>
            )}
            {isLoading &&
                <div className="col-12 d-flex justify-content-center align-items-center" style={{ height: 200 }}>
                    <Loading text={loadingText} />
                </div>
            }
        </div>
    );
}