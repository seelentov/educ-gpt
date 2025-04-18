'use client'

import { Dispatch, SetStateAction, useEffect, useState } from "react";
import Image from "next/image";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { me } from "@/core/api/auth/me";
import { updateUser } from "@/core/api/auth/update";
import { Loading } from "../ui/loading";
import { HOST_URL_PROD } from "@/core/api/api";
import { useRouter } from "next/navigation";
import { showToast } from "../utils/toast";

export default function EditProfileForm() {
    const [token] = useLocalStorage("token", "")

    const [name, setName] = useState<string>("")
    const [avatarFile, setAvatarFile] = useState<File | null>(null)

    const [avatarUrl, setAvatarUrl] = useState<string>("")

    const [errors, setErrors] = useState<{ [key: string]: string } | null>(null);
    const [loading, setLoading] = useState(false);
    const [loadingMe, setLoadingMe] = useState(false);

    const router = useRouter()

    useEffect(() => {
        if (token === "") {
            router.replace("/login")
        }
    })

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
                setName(data.name)

                if (data.avatar_url) {
                    setAvatarUrl(HOST_URL_PROD + data.avatar_url)
                }
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
                showToast("success", "Данные успешно обновлены")
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
                            </div>
                        </div>
                        <div className="col-12 d-flex justify-content-center align-items-center flex-column gap-1">
                            <button type="submit" className="btn btn-primary btn-block mb-4" disabled={loading}>
                                {loading ? "Отправка..." : "Сохранить"}
                            </button>
                            <p>{errors?.authorization}</p>
                        </div>
                    </form>
                </>}
        </>
    );
}
