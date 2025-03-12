import { Redirector } from "@/components/utils/redirector"
import { activate } from "@/core/api/auth/activate"

type Params = Promise<{ key: string }>

export default async function Activate(props: { params: Params }) {
    const params = await props.params;
    const data = await activate(params.key)

    return <div className="container mt-5">
        <p>{data?.status || data?.error}</p>
        <hr />
        {!data?.error && <Redirector timeout={5000} link="/login" />}
    </div>
}