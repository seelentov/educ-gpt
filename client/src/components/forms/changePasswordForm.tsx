'use client'

import { useState } from 'react';
import Link from 'next/link';
import { useLocalStorage } from '@/core/hooks/useLocalStorage';

export function ChangePasswordForm() {
    const [oldPassword, setOldPassword] = useState('');
    const [password, setPassword] = useState('');
    const [passwordConf, setPasswordConf] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false)

    const [token, settoken] = useLocalStorage("token", "")

    const handleSubmit = async (e: any) => {
        setLoading(true)
        setError('')
        e.preventDefault();

        try {

        } catch (error) {
            console.log(error)
            setError('Неверный пароль');
        }
        finally {
            setLoading(false)
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div className="d-flex flex-column mb-8">
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass_old">Пароль</label>
                    <input
                        className="input"
                        type="password"
                        id="fpass_old"
                        name="password_old"
                        placeholder="Пароль"
                        value={oldPassword}
                        onChange={(e) => setOldPassword(e.target.value)}
                    />
                </div>
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass">Новый пароль</label>
                    <input
                        className="input"
                        type="password"
                        id="fpass"
                        name="password"
                        placeholder="Новый пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass_new">Подтверждение пароля</label>
                    <input
                        className="input"
                        type="password"
                        id="fpass_new"
                        name="password_new"
                        placeholder="Подтверждение пароля"
                        value={passwordConf}
                        onChange={(e) => setPasswordConf(e.target.value)}
                    />
                </div>
            </div>

            {error && <div className="text-danger mb-4">{error}</div>}

            <hr />

            <div className="d-flex flex-wrap-wrap col-12 justify-content-between align-items-center">
                <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
                    {loading ? "В процессе..." : "Изменить пароль"}
                </button>
            </div>
        </form>
    );
}