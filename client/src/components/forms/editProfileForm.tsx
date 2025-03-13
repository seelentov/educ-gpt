'use client'

import { HOST_URL } from "@/core/api/api";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import Image from "next/image";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { me } from "@/core/api/auth/me";
import { updateUser } from "@/core/api/auth/update";
import { Loading } from "../ui/loading";

export default function EditProfileForm() {
    const [token] = useLocalStorage("token", "")

    const [name, setName] = useState<string>("")
    const [number, setNumber] = useState<string>("")
    const [avatarFile, setAvatarFile] = useState<File | null>(null)
    const [chatGptModel, setChatGptModel] = useState<string>("")
    const [chatGptToken, setChatGptToken] = useState<string>("")

    const [avatarUrl, setAvatarUrl] = useState<string>("")

    const [errors, setErrors] = useState<{ [key: string]: string } | null>(null);
    const [loading, setLoading] = useState(false);
    const [loadingMe, setLoadingMe] = useState(false);

    const [initState, setInitState] = useState<User | null>(null)

    const isDisabled =
        initState === null ||
        loading ||
        (
            initState.avatar_url === avatarUrl &&
            initState.chat_gpt_model === chatGptModel &&
            initState.name === name &&
            initState.number === number &&
            initState.chat_gpt_token === initState.chat_gpt_token
        )

    const refetchUser = async () => {
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
                setInitState(data)
                setName(data.name)
                setNumber(data.number)

                if (data.avatar_url) {
                    setAvatarUrl(data.avatar_url)
                }

                setChatGptModel(data.chat_gpt_model)
            }
        } catch (error) {
            console.error(error)
            setErrors({ authorization: JSON.stringify(error) });
        }
        finally {
            setLoadingMe(false)
        }
    }

    useEffect(() => {
        refetchUser()
    }, [])

    const updateData = async (e: any) => {
        if (loading) {
            return
        }

        try {
            setLoading(true)
            setErrors(null)
            e.preventDefault();

            const formData = new FormData()

            if (name) {
                formData.append("name", name)
            }
            if (number) {
                formData.append("number", number)
            }
            if (chatGptToken) {
                formData.append("chat_gpt_token", chatGptToken)
            }
            if (chatGptModel) {
                formData.append("chat_gpt_model", chatGptModel)
            }
            if (name) {
                formData.append("name", name)
            }
            if (avatarFile) {
                formData.append("avatar_file", avatarFile)
            }

            const data = await updateUser(formData, token);

            if (data?.error) {
                if (typeof data?.error === "string") {
                    setErrors({ authorization: data.error })
                } else {
                    setErrors({ ...data.error, authorization: data.error.avatar_file })
                }
            }
            else {
                alert("Данные успешно обновлены")
                await refetchUser()
            }
        } catch (error) {
            console.error(error)
            setErrors({ authorization: JSON.stringify(error) });
        }
        finally {
            setLoading(false)
        }
    }

    function handleChange<T>(dispatch: Dispatch<SetStateAction<T>>, value: T) {
        dispatch(value)
    }

    const uploadImage = (e: React.ChangeEvent<HTMLInputElement>) => {
        const targetFiles = e.target.files
        if (targetFiles === null) {
            return
        }
        const files: File[] = Array.from(targetFiles)
        if (files.length < 1) {
            return
        }
        const file = files[0]

        setAvatarFile(file)

        const imageUrl = URL.createObjectURL(file);

        setAvatarUrl(imageUrl)
    }

    console.log(avatarUrl)

    return (
        <>
            {loadingMe
                ? <div className="col-12 d-flex justify-content-center align-items-center" style={{ height: 200 }}>
                    <Loading />
                </div>
                : <>
                    <form onSubmit={updateData} className="col-12 d-flex align-items-center flex-wrap">
                        <div className="col-12 col-md-6 avatar">
                            <label className="avatar_edit-btn">
                                <Image style={{ objectFit: 'cover' }} src={"/icons/edit.svg"} alt="" width={60} height={60} />
                                <input className="avatar_edit-btn" type="file" onChange={uploadImage} accept="image/*" hidden></input>
                            </label>
                            {
                                avatarUrl && avatarUrl !== ""
                                    ? <Image style={{ objectFit: 'cover' }} src={avatarUrl} alt="" width={300} height={300} className="rounded-circle" />
                                    : <Image style={{ objectFit: 'cover' }} src={"/misc/empty_avatar.jpg"} alt="" width={300} height={300} className="rounded-circle" />
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
                                        onChange={(e) => handleChange(setName, e.target.value)}
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
                                        onChange={(e) => handleChange(setNumber, e.target.value)}
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
                                        onChange={(e) => handleChange(setChatGptToken, e.target.value)}
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
                                        value={chatGptModel}
                                        onChange={(e) => handleChange(setChatGptModel, e.target.value)}
                                    />
                                    <p className="text-danger">{errors?.chat_gpt_model}</p>
                                </div>
                            </div>
                        </div>
                        <div className="col-12 d-flex justify-content-center align-items-center flex-column gap-1">
                            <button type="submit" className="btn btn-primary btn-block mb-4" disabled={isDisabled}>
                                {loading ? "Отправка..." : "Сохранить"}
                            </button>
                            <p>{errors?.authorization}</p>
                        </div>
                    </form>
                </>}
        </>
    );
}
