import { ChangeEmailForm } from "@/components/forms/changeEmailForm";
import { ChangePasswordForm } from "@/components/forms/changePasswordForm";
import EditProfileForm from "@/components/forms/editProfileForm";

export default function Profile() {
    return (
        <div className="container mt-5">
            <div className="text">
                <h1>Настройки профиля</h1>
            </div>
            <EditProfileForm />
            <hr />
            <div className="text">
                <h2>Сменить пароль</h2>
            </div>
            <ChangePasswordForm />
            <hr />
            <div className="text">
                <h2>Сменить почту</h2>
            </div>
            <ChangeEmailForm />
        </div>
    );
}
