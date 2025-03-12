'use client'

import { useEffect, useMemo, useState } from "react";
import { remark } from 'remark';
import html from 'remark-html';
import rehypeHihglight from 'rehype-highlight'
import { useParams, useRouter } from "next/navigation";
import CodeMirror from "@uiw/react-codemirror";
import { getTheme } from "@/core/api/roadmap/theme";
import { useLocalStorage } from "@/core/hooks/useLocalStorage";
import { Loading } from "@/components/ui/loading";
import { resolve } from "@/core/api/roadmap/resolve";

export default function ThemePage() {
    const [part, setPath] = useState<"theory" | "tasks">("theory")

    const [content, setContent] = useState<string>("")

    const [tasks, setTasks] = useState<{ task: string, isDone: boolean, dialog: string[], id: number }[]>([])
    const [activeTask, setActiveTask] = useState<number>(0)

    const [code, setCode] = useState<string>("")
    const [consoleText, setConsoleText] = useState<string>()
    const router = useRouter()
    const { topic, theme } = useParams()

    const [token] = useLocalStorage("token", "")
    const [globalLoading, setGlobalLoading] = useState<boolean>(true)
    const [problemsLoading, setProblemsLoading] = useState<boolean>(true)
    const [checkLoading, setCheckLoading] = useState<boolean>(true)

    useEffect(() => {
        (async () => {
            if (!topic || !theme || !token) {
                return
            }
            setGlobalLoading(true)
            const data = await getTheme(topic as string, theme as string, token)

            if (data?.error) {
                console.error(data.error)
                router.refresh()
                return
            }

            const processedContent = await remark()
                .use(html)
                .use(rehypeHihglight)
                .process(data.text);
            setContent(processedContent.toString())

            for (let i = 0; i < data.problems.length; i++) {
                const processedProblem = await remark()
                    .use(html)
                    .use(rehypeHihglight)
                    .process(data.problems[i].question);

                data.problems[i].question = processedProblem.toString()
            }

            setTasks(data.problems.map((p: Problem) => {
                return {
                    id: p.id,
                    task: p.question,
                    isDone: false,
                    dialog: []
                }
            }))
            setGlobalLoading(false)
        })()
    }, [topic, theme, token])

    const checkAnswer = async () => {
        if (checkLoading || tasks[activeTask].isDone) {
            return
        }

        setCheckLoading(true)

        const data = await resolve(tasks[activeTask].id, code, token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            return
        }

        if (data?.message && data?.status) {
            setTasks(ts => ([
                ...ts.filter(t => t.id !== tasks[activeTask].id),
                {
                    ...tasks[activeTask],
                    dialog: [...tasks[activeTask].dialog, data.message],
                    isDone: data.status
                }
            ]))
        } else {
            console.error(data)
            alert(JSON.stringify(data))
        }

        setCheckLoading(false)
    }

    const loadMoreProblems = async () => {
        if (problemsLoading) {
            return
        }

        const data = await getTheme(topic as string, theme as string, token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            return
        }

        setTasks(ts => ([
            ...ts,
            ...data.map((p: Problem) => {
                return {
                    task: p.question,
                    isDone: false,
                    dialog: []
                }
            })
        ]))

        setProblemsLoading(true)
        setProblemsLoading(false)
    }

    const canNextTask = useMemo(() => tasks.length > activeTask + 1, [tasks, activeTask])

    const nextTask = () => {
        if (!canNextTask || checkLoading) {
            return
        }

        setActiveTask(p => p + 1)
    }

    const canPrevTask = useMemo(() => activeTask > 0, [tasks, activeTask])

    const prevTask = () => {
        if (!canPrevTask || checkLoading) {
            return
        }

        setActiveTask(p => p - 1)
    }

    return (
        <>
            {!globalLoading
                ? <div className="col12 p-3 border d-flex gap-3 pt-5 position-relative" style={{ height: '100vh' }}>
                    <div className="position-absolute fixed-top p-2 d-flex" style={{ height: '48px' }}>
                        <button onClick={() => router.back()} type="button" className={`btn btn-outline-primary btn-sm me-2`}>Назад</button>
                        <button onClick={() => setPath('theory')} style={{ marginLeft: 'auto' }} type="button" className={`btn ${part === 'theory' ? "btn-primary" : "btn-outline-primary"} btn-sm me-2`}>Теория</button>
                        <button onClick={() => setPath('tasks')} type="button" className={`btn ${part === 'tasks' ? "btn-primary" : "btn-outline-primary"} btn-sm me-2`}>Задачи</button>

                    </div>
                    <div className="col-6 border rounded position-relative" style={{ height: '100%' }}>
                        <button type="button" className="btn btn-success btn-sm position-absolute" style={{ top: '10px', right: '10px', zIndex: 999 }}>{"▶"}</button>
                        <div className="border" style={{ height: 'calc(100% - 150px)', overflow: 'scroll' }}>
                            <CodeMirror
                                onChange={setCode}
                                value={code}
                                height="100%"
                            />
                        </div>
                        <div className="border p-2 " style={{ height: '150px', overflow: 'scroll' }}>
                            {consoleText}
                        </div>
                    </div>
                    <div className="col-6 border rounded p-2" style={{ height: '100%', overflow: 'scroll' }}>
                        {part === 'tasks' && tasks && tasks.length > 0
                            ? <>
                                <div className="text p-2 border rounded" style={{ height: 'calc(100% - 50px)', overflow: 'scroll' }}>
                                    <div dangerouslySetInnerHTML={{ __html: tasks[activeTask].task }} />
                                    {tasks[activeTask].dialog.map((d) => <p><strong>AI:</strong> {d}</p>)}
                                </div>
                                <div className="p-2 border rounded d-flex gap-1 align-items-center" style={{ height: '50px' }}>
                                    <button onClick={() => prevTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canPrevTask}>{"<<"}</button>
                                    <button disabled={checkLoading || tasks[activeTask].isDone} onClick={() => checkAnswer()} type="button" className="btn btn-outline-success btn-sm">Проверить решение</button>
                                    <button onClick={() => nextTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canNextTask}>{">>"}</button>
                                    <p className="test-center mx-auto mb-0">{activeTask + 1} / {tasks.length}</p>
                                    <button onClick={() => loadMoreProblems()} disabled={problemsLoading} type="button" className="btn btn-outline-warning btn-sm" style={{ marginLeft: 'auto' }}>Загрузить еще...</button>
                                </div>
                            </>
                            : <>
                                <div className="text" dangerouslySetInnerHTML={{ __html: content }} />
                            </>}

                    </div>
                </div >
                : <div className="d-flex justify-content-center align-items-center" style={{ height: "100vh" }}>
                    <Loading text={"AI подбирает текст и информацию, учитывая ваши успехи"} />
                </div>

            }
        </>

    )

}
