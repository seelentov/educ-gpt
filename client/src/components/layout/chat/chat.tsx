'use client'

import { useEffect, useMemo, useState, useRef } from "react"
import { usePathname } from 'next/navigation'

import Image from "next/image";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { deleteDialog } from "@/core/api/dialogs/delete_dialog";
import { getDialogs } from "@/core/api/dialogs/get_dialogs";
import { getDialog } from "@/core/api/dialogs/get_dialog";
import { Loading } from "@/components/ui/loading";
import { createDialog } from "@/core/api/dialogs/create_dialog";
import { sendMessage } from "@/core/api/dialogs/send_message";

export function Chat() {
    const [isOpen, setIsOpen] = useState<boolean>(false)
    const [isThrowLoading, setIsThrowLoading] = useState<boolean>(false)

    const [input, setInput] = useState<string>("")

    const [activeDialogIndex, setActiveDialogIndex] = useState<number>(0)

    const [messages, setMessages] = useState<DialogItem[]>([])
    const [isMessagesLoading, setIsMessagesLoading] = useState<boolean>(true)

    const [dialogs, setDialogs] = useState<Dialog[]>([])
    const [isDialogsLoading, setIsDialogsLoading] = useState<boolean>(true)

    const [createNewDialogLoading, setCreateNewDialogLoading] = useState<boolean>(false)

    const activeDialog: Dialog = useMemo(() => dialogs[activeDialogIndex], [dialogs, activeDialogIndex])

    const deleteDialogDisabled = useMemo(() => dialogs.length < 2, [dialogs])

    const [token, setToken] = useLocalStorage("token", "")

    const pathname = usePathname()

    const messagesEndRef = useRef<HTMLDivElement>(null)

    const scrollMessagesToEnd = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    useEffect(() => {
        scrollMessagesToEnd()
    }, [messages])

    const fetchDialogs = async () => {
        setIsDialogsLoading(true)

        const res = await getDialogs(token)
        if (res?.error) {
            const err = res.error;
            console.error(err)
            alert(err)
        }

        setDialogs(res)
        setMessages(res.dialog_items)

        setIsDialogsLoading(false)
    }

    const fetchDialog = async () => {
        setIsMessagesLoading(true)

        const res = await getDialog(activeDialog.id, token)
        if (res?.error) {
            const err = res.error;
            console.error(err)
            alert(err)
        }

        setMessages(res.dialog_items)
        scrollMessagesToEnd()

        setIsMessagesLoading(false)
    }

    const throwMessageHandle = async (e: any) => {
        e.preventDefault()

        if (isThrowLoading) {
            return
        }

        setInput("")

        setMessages(p => ([...p,
        {
            id: 9999999 + p.length,
            text: input,
            is_user: true
        }]))

        setIsThrowLoading(true)

        const res = await sendMessage(activeDialog.id, input, token)
        if (res?.error) {
            const err = res.error;
            console.error(err)
            alert(err)
        }

        setMessages(p => ([...p, res]))

        scrollMessagesToEnd()

        setIsThrowLoading(false)
    }

    const deleteDialogHandle = async (id: number) => {
        if (deleteDialogDisabled) {
            return
        }

        setDialogs((p) => ([...p.filter(d => d.id != id)]))

        if (id === activeDialog.id) {
            setActiveDialogIndex(0)
        }

        const res = await deleteDialog(id, token)

        if (res?.error) {
            const err = res.error;
            console.error(err)
            alert(err)
        }
    }

    const createDialogHandle = async () => {
        console.log(createNewDialogLoading)

        if (createNewDialogLoading) {
            return
        }

        setCreateNewDialogLoading(true)

        const res = await createDialog(token)
        if (res?.error) {
            const err = res.error;
            console.error(err)
            alert(err)
        }

        setDialogs(p => ([...p, res]))

        setCreateNewDialogLoading(false)
    }

    const setDialog = async (i: number) => {
        if (isMessagesLoading) {
            return
        }

        setIsMessagesLoading(true)

        setActiveDialogIndex(i)
        await fetchDialog()

        setIsMessagesLoading(false)
    }

    useEffect(() => {
        fetchDialogs()
        setActiveDialogIndex(0)
    }, [token])

    useEffect(() => {
        if (dialogs.length > 0) {
            fetchDialog()
        }
    }, [activeDialogIndex, token, dialogs])

    useEffect(() => {
        try {
            const value = window.localStorage.getItem("token")

            const v = value ? JSON.parse(value) : ""

            if (v !== "") {
                setToken(v)
            }

        } catch (error) {
            console.error(error)
        }

    }, [pathname])

    return (
        <div className="position-fixed" style={{ bottom: 20, right: 20 }}>
            {
                !isOpen
                    ? <button onClick={() => setIsOpen(true)} type="button" className="btn btn-primary" style={{ zIndex: 999 }}>
                        <Image src="/icons/chat.svg" alt="" width={20} height={20} />
                    </button>
                    : <div className="container-fluid">
                        <div className="row justify-content-center">
                            <div className="col-6 p-0" style={{ width: '300px', height: '400px', border: '1px solid #ccc', borderRadius: '5px', overflow: 'hidden' }}>
                                <div className="d-flex flex-column h-100">
                                    <div className="p-2 d-flex gap-1 align-items-center" style={{ backgroundColor: '#e9ecef' }}>
                                        <button onClick={() => setIsOpen(false)} className="btn btn-secondary btn-sm">
                                            ×
                                        </button>
                                        {isDialogsLoading
                                            ? <Loading min color="white" />
                                            : <>

                                                <div className="d-flex h-100 gap-1 w-100" style={{ overflowX: 'scroll', height: '30px !important' }}>
                                                    {dialogs.map((d, i) =>
                                                        <div
                                                            className={`d-flex rounded-1 overflow-hidden`}
                                                            key={d.id}
                                                            style={{ padding: 0, minWidth: 90 }}
                                                        >
                                                            <button
                                                                onClick={() => deleteDialogHandle(d.id)}
                                                                className={`btn btn-danger btn-sm rounded-0`}
                                                                disabled={deleteDialogDisabled}
                                                            >
                                                                ×
                                                            </button>
                                                            <button
                                                                onClick={() => setDialog(i)}
                                                                className={`btn btn-primary btn-sm rounded-0 border-0 w-100`}
                                                                disabled={activeDialogIndex === i}
                                                            >
                                                                {d.dialog_items.length > 0 ? d.dialog_items[0].text.slice(0, 7) + "..." : ""}
                                                            </button>
                                                        </div>
                                                    )}
                                                    <button onClick={createDialogHandle} className="btn btn-success btn-sm">
                                                        {createNewDialogLoading ? <Loading min color="white" /> : "+"}
                                                    </button>
                                                </div>
                                            </>}
                                    </div>
                                    <div className="flex-grow-1 p-2" style={{ overflowY: 'auto', backgroundColor: 'white' }}>
                                        {isMessagesLoading
                                            ? <Loading min />
                                            : messages.map(m =>
                                                <div className="mb-2" key={m.id}>
                                                    <div className="bg-light p-2 rounded">
                                                        <strong>{m.is_user ? "User" : "AI"}:</strong> {m.text}
                                                    </div>
                                                </div>
                                            )}
                                        {isThrowLoading
                                            && <div className="mb-2">
                                                <div className="bg-light p-2 rounded">
                                                    <strong>AI:</strong> печатаю...
                                                </div>
                                            </div>}
                                        <div ref={messagesEndRef} />
                                    </div>
                                    <div className="p-2" style={{ backgroundColor: '#e9ecef' }}>
                                        <form onSubmit={throwMessageHandle}>
                                            <input
                                                type="text"
                                                className="form-control"
                                                placeholder="Введите сообщение..."
                                                value={input}
                                                onChange={(e) => setInput(e.target.value)}
                                                disabled={isThrowLoading}
                                            />
                                        </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
            }
        </div >
    )
}