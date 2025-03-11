'use client'

import { useState } from 'react';
import { resetPasswordActivate } from '@/core/api/auth/reset_password_activate';
import { useParams } from 'next/navigation';

export function ResetPasswordForm() {
    const [password, setPassword] = useState('');
    const [passwordConf, setPasswordConf] = useState('');

    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState<string>("")

    const { key, userId } = useParams()

    const handleSubmit = async (e: any) => {
        setLoading(true)
        setError('')
        e.preventDefault();

        if (password != passwordConf) {
            setError("Пароль и его подтверждение должны быть идентичны")
            setLoading(false)
            return
        }

        if (!key || !userId) {
            return
        }

        try {
            const data = await resetPasswordActivate(key as string, userId as string, password)
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setError(data?.error)
                } else {
                    setError(data.error[Object.keys(data.error)[0]])
                }
            }
            else {
                setMessage(data?.status)
            }
        } catch (error) {
            console.error(error)
            setError(JSON.stringify(error));
        }
        finally {
            setLoading(false)
            setPassword("")
            setPasswordConf("")
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            {message !== ""
                ? <p>{message}</p>
                : <>
                    <div className="d-flex flex-column mb-2 gap-3">
                        <label htmlFor="fpass">Пароль</label>
                        <input
                            className="input"
                            type="password"
                            id="fpass"
                            name="password"
                            placeholder="Пароль"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <div className="d-flex flex-column mb-2 gap-3">
                        <label htmlFor="fpassconf">Подтверждение пароля</label>
                        <input
                            className="input"
                            type="password"
                            id="fpassconf"
                            name="passwordconf"
                            placeholder="Подтверждение пароля"
                            value={passwordConf}
                            onChange={(e) => setPasswordConf(e.target.value)}
                        />
                    </div>

                    {error && <div className="text-danger mb-4">{error}</div>}

                    <div className="d-flex flex-wrap-wrap col-12 justify-content-between align-items-center">
                        <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
                            {loading ? "В процессе..." : "Сменить пароль"}
                        </button>
                    </div>
                </>}
        </form>
    );
}