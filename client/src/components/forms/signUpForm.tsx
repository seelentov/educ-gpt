'use client'

import { signup } from "@/core/api/auth/signup";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from 'next/navigation'
import { useLocalStorage } from "@/core/hooks/useLocalStorage";

export function SignUpForm() {

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [number, setNumber] = useState('');
    const [password, setPassword] = useState('');
    const [passwordConf, setPasswordConf] = useState('');
    const [chatGptToken, setChatGptToken] = useState('');

    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState<{ [key: string]: string } | null>(null);

    const [success, setSuccess] = useState(false)

    const [token] = useLocalStorage("token", "")

    const router = useRouter()

    useEffect(() => {
        if (token !== "") {
            router.push("/topics")
        }
    })

    const handleSubmit = async (e: any) => {
        if (loading) {
            return
        }

        setLoading(true)
        setErrors(null)
        e.preventDefault();

        if (password != passwordConf) {
            setErrors({ passwordConf: "Пароль и его подтверждение должны быть идентичны" })
            setLoading(false)
            return
        }

        try {
            const data = await signup(
                name,
                email,
                number,
                password,
                chatGptToken
            );
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setErrors({ authorization: data.error })
                } else {
                    setErrors(data?.error)
                }
            }
            else {
                setSuccess(true)
            }
        } catch (error) {
            console.error(error)
            setErrors({ authorization: JSON.stringify(error) });
        }
        finally {
            setLoading(false)
            setPassword("")
            setPasswordConf("")
        }
    };

    return (
        <>
            {success
                ? <p>Регистрация прошла успешно. Подтвердите аккаунт через письмо на почте {email}</p>
                : <form onSubmit={handleSubmit}>
                    <div className="d-flex flex-column mb-8">
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="fname">Имя</label>
                            <input
                                className={`input ${errors?.name && 'err'}`}
                                type="text"
                                id="fname"
                                name="name"
                                placeholder="Имя"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                            <p className="text-danger">{errors?.name}</p>
                        </div>
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="femail">E-mail</label>
                            <input
                                type="email"
                                className={`input ${errors?.email && 'err'}`}
                                id="femail"
                                name="email"
                                placeholder="E-mail"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                            <p className="text-danger">{errors?.email}</p>
                        </div>
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="fnumber">Номер телефона</label>
                            <input
                                type="tel"
                                className={`input ${errors?.number && 'err'}`}
                                id="fnumber"
                                name="number"
                                placeholder="Телефон"
                                value={number}
                                onChange={(e) => setNumber(e.target.value)}
                            />
                            <p className="text-danger">{errors?.number}</p>
                        </div>
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="fpass">Пароль</label>
                            <input
                                className={`input ${(errors?.password || errors?.passwordConf) && 'err'}`}
                                type="password"
                                id="fpass"
                                name="password"
                                placeholder="Пароль"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                            <p className="text-danger">{errors?.password}</p>
                        </div>
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="fpassconf">Подтверждение пароля</label>
                            <input
                                className={`input ${errors?.password || errors?.passwordConf && 'err'}`}
                                type="password"
                                id="fpassconf"
                                name="passwordconf"
                                placeholder="Подтверждение пароля"
                                value={passwordConf}
                                onChange={(e) => setPasswordConf(e.target.value)}
                            />
                            <p className="text-danger">{errors?.passwordConf}</p>
                        </div>
                        <div className="d-flex flex-column mb-2 gap-3">
                            <label htmlFor="fgpttoken">Ключ Chat-GPT API</label>
                            <input
                                className={`input ${errors?.chat_gpt_token && 'err'}`}
                                type="password"
                                id="fgpttoken"
                                name="gpttoken"
                                placeholder="Ключ Chat-GPT API"
                                value={chatGptToken}
                                onChange={(e) => setChatGptToken(e.target.value)}
                            />
                            <p>Получить ключ на <a className="link-primary" href="https://platform.openai.com/api-keys" target="_blank">openai.com</a></p>
                            <p className="text-danger">{errors?.chat_gpt_token}</p>
                        </div>
                    </div>

                    <hr />

                    <p className="text-danger">{errors?.authorization}</p>

                    <button type="submit" className="btn btn-primary btn-block mb-4" disabled={loading}>
                        {loading ? "Регистрация..." : "Присоединиться"}
                    </button>

                    <hr />

                    <div className="text-center">
                        <p>Есть аккаунт? <Link className="link-primary" href="/login">Войти</Link></p>
                    </div>
                </form>}
        </>
    )
}