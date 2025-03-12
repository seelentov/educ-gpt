import { changeEmailActivate } from "@/core/api/auth/change_email_activate";

type Params = Promise<{ key: string, userId: string }>

export default async function ChangeEmailActivatePage(props: { params: Params }) {
    const params = await props.params;
    const data = await changeEmailActivate(params.key, params.userId)

    return <div className="container mt-5">
        <p>{data?.status || data?.error}</p>
        <hr />
    </div>
}