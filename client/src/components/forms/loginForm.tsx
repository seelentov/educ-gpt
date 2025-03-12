'use client'

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { login } from '@/core/api/auth/login';
import { useLocalStorage } from '@/core/hooks/useLocalStorage';
import { useRouter } from 'next/navigation'

export function LoginForm() {
    const [credential, setCredential] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [token, setToken] = useLocalStorage("token", "")
    const [, setRefreshToken] = useLocalStorage("refresh_token", "")
    const [loading, setLoading] = useState(false)

    const router = useRouter()

    useEffect(() => {
        if (token != "") {
            router.replace("/topics")
        }
    })


    const handleSubmit = async (e: any) => {
        if (loading) {
            return
        }

        setLoading(true)
        setError('')
        e.preventDefault();

        try {
            const data = await login(credential, password);
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setError(data?.error)
                } else {
                    setError(data.error[Object.keys(data.error)[0]])
                }
            }
            else {
                if (data?.token) {
                    setToken(data.token)
                }
                if (data?.refresh_token) {
                    setRefreshToken(data.refresh_token)
                }

                router.replace("/topics");
            }


        } catch (error) {
            console.error(error)
            setError('Неверный логин или пароль');
        }
        finally {
            setLoading(false)
            setPassword("")
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div className="d-flex flex-column mb-8">
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fname">Логин</label>
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
                <div className="d-flex flex-column mb-4 gap-3">
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
            </div>

            {error && <div className="text-danger mb-4">{error}</div>}

            <hr />

            <div className="d-flex flex-wrap-wrap col-12 justify-content-between align-items-center">
                <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
                    {loading ? "Вход..." : "Войти"}
                </button>
                <Link href="/reset" className='link-primary'>Забыли пароль?</Link>
            </div>

            <hr />

            <div className="text-center">
                <p>Нет аккаунта? <Link className="link-primary" href="/signup">Регистрация</Link></p>
            </div>
        </form>
    );
}