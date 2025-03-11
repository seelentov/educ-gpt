'use client'

import { useState } from 'react';
import Link from 'next/link';
import { resetPassword } from '@/core/api/auth/reset_password';

export function ResetPasswordTaskForm() {
    const [credential, setCredential] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState<string>("")


    const handleSubmit = async (e: any) => {

        setLoading(true)
        setError('')
        e.preventDefault();

        try {
            const data = await resetPassword(credential)
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setError(data?.error)
                } else {
                    setError(data.error[Object.keys(data.error)[0]])
                }
            }
            else {
                setMessage(data?.message)
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
            {message !== ""
                ? <p>{message}</p>
                : <>
                    <div className="d-flex flex-column mb-8">
                        <div className="d-flex flex-column mb-4 gap-3">
                            <label htmlFor="fname">Email / Телефон / Имя</label>
                            <input
                                className="input"
                                type="text"
                                id="fname"
                                name="login"
                                placeholder="E-mail / Телефон / Имя"
                                value={credential}
                                onChange={(e) => setCredential(e.target.value)}
                            />
                        </div>
                    </div>

                    {error && <div className="text-danger mb-4">{error}</div>}

                    <div className="d-flex flex-wrap-wrap col-12 justify-content-between align-items-center">
                        <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
                            {loading ? "Вход..." : "Войти"}
                        </button>
                    </div>
                </>}
        </form>
    );
}