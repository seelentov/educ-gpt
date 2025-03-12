import { SignUpForm } from "@/components/forms/signUpForm";

export default function Signup() {
    return (
        <div className="container">
            <div className="d-flex justify-content-center mt-5">
                <div className="col-5">
                    <SignUpForm />
                </div>
            </div>
        </div>
    );
}
