'use client'

import { HOST_URL } from "@/core/api/api";
import { useEffect, useState } from "react";
import Image from "next/image";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { me } from "@/core/api/auth/me";
import { updateUser } from "@/core/api/auth/update";
import { Loading } from "../ui/loading";

export default function EditProfileForm() {
    const [token, _] = useLocalStorage("token", "")

    const [name, setName] = useState<string>("")
    const [number, setNumber] = useState<string>("")
    const [avatarUrl, setAvatarUrl] = useState<string>("")
    const [chatGptModel, setChatGptModel] = useState<string>("")
    const [chatGptToken, setChatGptToken] = useState<string>("")

    const [errors, setErrors] = useState<{ [key: string]: string } | null>(null);
    const [loading, setLoading] = useState(false);
    const [loadingMe, setLoadingMe] = useState(false);


    useEffect(() => {
        (async () => {
            setLoadingMe(true)

            try {
                const data = await me(token)

                if (data?.error) {
                    if (typeof data?.error === "string") {
                        setErrors({ authorization: data.error })
                    } else {
                        setErrors(data?.error)
                    }
                }
                else {
                    setName(data.name)
                    setNumber(data.number)
                    setAvatarUrl(data.avatar_url)
                    setChatGptModel(data.chat_gpt_model)
                    setChatGptToken(data.chat_gpt_token)
                }
            } catch (error) {
                console.log(error)
                setErrors({ authorization: JSON.stringify(error) });
            }
            finally {
                setLoadingMe(false)
            }
        })()
    }, [token])

    const updateData = async (e: any) => {
        try {
            setLoading(true)
            setErrors(null)
            e.preventDefault();

            const data = await updateUser(
                name,
                number,
                chatGptToken,
                chatGptModel,
                avatarUrl,
                token
            );
            if (data?.error) {
                if (typeof data?.error === "string") {
                    setErrors({ authorization: data.error })
                } else {
                    setErrors(data?.error)
                }
            }
            else {
                alert("Данные успешно обновлены")
            }
        } catch (error) {
            console.log(error)
            setErrors({ authorization: JSON.stringify(error) });
        }
        finally {
            setLoading(false)
        }
    }

    const uploadImage = () => {

    }

    return (
        <>
            {loadingMe
                ? <div className="col-12 d-flex justify-content-center align-items-center" style={{ height: 200 }}>
                    <Loading />
                </div>
                : <>
                    <form onSubmit={updateData} className="col-12">
                        <div className="col-12 col-md-6 avatar">
                            <div className="avatar_edit-btn" onClick={uploadImage}>
                                <Image src={"/icons/edit"} alt="" width={30} height={30} className="rounded-circle" />
                            </div>
                            {
                                avatarUrl && avatarUrl !== ""
                                    ? <Image src={HOST_URL + "/storage/" + avatarUrl} alt="" width={300} height={300} className="rounded-circle" />
                                    : <Image src={"/misc/empty_avatar.jpg"} alt="" width={300} height={300} className="rounded-circle" />
                            }
                        </div>
                        <div className="col-12 col-md-6">
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

                                <div className="d-flex flex-column mb-2 gap-3">
                                    <label htmlFor="fgptmodel">Модель Chat-GPT</label>
                                    <input
                                        className={`input ${errors?.chat_gpt_model && 'err'}`}
                                        type="password"
                                        id="fgptmodel"
                                        name="gptmodel"
                                        placeholder="Модель Chat-GPT"
                                        value={chatGptToken}
                                        onChange={(e) => setChatGptToken(e.target.value)}
                                    />
                                    <p className="text-danger">{errors?.chat_gpt_model}</p>
                                </div>
                            </div>
                        </div>
                    </form>
                    <div className="col-12">
                        <button type="submit" className="btn btn-primary btn-block mb-4" disabled={loading}>
                            {loading ? "Отправка..." : "Изменить данные"}
                        </button>
                    </div>
                </>}
        </>
    );
}
