'use client'

import { changePassword } from '@/core/api/auth/change_password';
import { useLocalStorage } from '@/core/hooks/useLocalStorage';
import { useState } from 'react';

export function ChangePasswordForm() {
    const [oldPassword, setOldPassword] = useState('');
    const [password, setPassword] = useState('');
    const [passwordConf, setPasswordConf] = useState('');
    const [errors, setErrors] = useState<{ [key: string]: string } | null>(null);
    const [loading, setLoading] = useState(false)

    const [token] = useLocalStorage("token", "")

    const disabled = !oldPassword && !password && !passwordConf

    const handleSubmit = async (e: any) => {
        setLoading(true)
        setErrors(null)
        e.preventDefault();

        if (password != passwordConf) {
            setErrors({ passwordConf: "Пароль и его подтверждение должны быть идентичны" })
            setLoading(false)
            return
        }

        try {
            const data = await changePassword(oldPassword, password, token);
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setErrors({ authorization: data.error })
                } else {
                    setErrors(data?.error)
                }
            }
            else {
                alert("Пароль успешно изменен")
            }

            console.log(data)
        } catch (error) {
            console.error(JSON.stringify(error))
            setErrors({ authorization: JSON.stringify(error) });
        }
        finally {
            setOldPassword("")
            setPassword("")
            setPasswordConf("")
            setLoading(false)
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div className="d-flex flex-column mb-8">
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass_old">Пароль</label>
                    <input
                        className={`input ${errors?.old_password && 'err'}`}
                        type="password"
                        id="fpass_old"
                        name="password_old"
                        placeholder="Пароль"
                        value={oldPassword}
                        onChange={(e) => setOldPassword(e.target.value)}
                    />
                </div>
                <p className="text-danger">{errors?.old_password}</p>
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass">Новый пароль</label>
                    <input
                        className={`input ${errors?.password && 'err'}`}
                        type="password"
                        id="fpass"
                        name="password"
                        placeholder="Новый пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                <p className="text-danger">{errors?.password}</p>
                <div className="d-flex flex-column mb-4 gap-3">
                    <label htmlFor="fpass_new">Подтверждение пароля</label>
                    <input
                        className={`input ${errors?.password && 'err'}`}
                        type="password"
                        id="fpass_new"
                        name="password_new"
                        placeholder="Подтверждение пароля"
                        value={passwordConf}
                        onChange={(e) => setPasswordConf(e.target.value)}
                    />
                </div>
                <p className="text-danger">{errors?.password}</p>
            </div>

            <div className="col-12 d-flex justify-content-center align-items-center flex-column gap-1">
                <button type="submit" className="btn btn-primary btn-block" disabled={loading || disabled}>
                    {loading ? "В процессе..." : "Изменить пароль"}
                </button>
                <br />
                <p className="text-danger">{errors?.authorization}</p>
            </div>
        </form>
    );
}