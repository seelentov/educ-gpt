'use client'

import { changeEmail } from "@/core/api/auth/change_email";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { useState } from "react";

export function ChangeEmailForm() {
    const [email, setEmail] = useState<string>("")

    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false)

    const [token] = useLocalStorage("token", "")

    const handleSubmit = async (e: any) => {
        if (loading) {
            return
        }

        setLoading(true)
        setError('')
        e.preventDefault();

        try {
            const data = await changeEmail(email, token);
            if (data?.error) {
                setError(data?.error)
            }
            else {
                alert(data?.message)
            }
        } catch (error) {
            console.error(error)
            setError(JSON.stringify(error));
        }
        finally {
            setLoading(false)
        }
    };


    return (
        <form onSubmit={handleSubmit}>
            <div className="d-flex flex-column flex-md-row flex-wrap gap-1 justify-content-center align-items-md-center align-items-initial">
                <input
                    className="input"
                    type="text"
                    id="femail"
                    name="email"
                    placeholder="E-mail"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
                    {loading ? "Отправка..." : "Изменить почту"}
                </button>
            </div>
            <p className="text-danger mt-2 text-center">{error}</p>
        </form>
    )
}