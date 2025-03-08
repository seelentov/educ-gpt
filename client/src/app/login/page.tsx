import { LoginForm } from "@/components/forms/loginForm";

export default function Login() {
    return (
        <div className="container">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-5">
                    <LoginForm />
                </div>
            </div>
        </div>
    );
}
