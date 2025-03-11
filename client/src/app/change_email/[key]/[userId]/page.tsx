import { changeEmailActivate } from "@/core/api/auth/change_email_activate";

interface IChangeEmailActivateParams {
    params: {
        key: string,
        userId: string
    }
}

export default async function ChangeEmailActivatePage({ params }: IChangeEmailActivateParams) {
    const data = await changeEmailActivate(params.key, params.userId)

    return <div className="container mt-5">
        <p>{data?.status || data?.error}</p>
        <hr />
    </div>
}