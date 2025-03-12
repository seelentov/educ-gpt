import { ResetPasswordForm } from "@/components/forms/resetPasswordForm"

export default async function ResetPasswordActivate() {

    return (
        <div className="conteiner mt-5">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-12 col-md-5">
                    <ResetPasswordForm />
                </div>
            </div>
        </div>
    )
}