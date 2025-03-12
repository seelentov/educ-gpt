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
import { compile } from "@/core/api/utils/compile";
import { getProblems } from "@/core/api/roadmap/problems";

export default function ThemePage() {
    const [part, setPath] = useState<"theory" | "tasks">("theory")

    const [content, setContent] = useState<string>("")

    const [tasks, setTasks] = useState<{ task: string, isDone: boolean, dialog: string[], id: number }[]>([])
    const [activeTask, setActiveTask] = useState<number>(0)

    const [code, setCode] = useState<string>("")
    const [consoleText, setConsoleText] = useState<string[]>([])
    const router = useRouter()
    const { topic, theme } = useParams()

    const [token] = useLocalStorage("token", "")
    const [globalLoading, setGlobalLoading] = useState<boolean>(true)
    const [problemsLoading, setProblemsLoading] = useState<boolean>(false)
    const [checkLoading, setCheckLoading] = useState<boolean>(false)
    const [compilationLoading, setCompilationLoading] = useState<boolean>(false)

    useEffect(() => {
        (async () => {
            if (!topic || !theme || !token) {
                return
            }
            setGlobalLoading(true)
            const data = await getTheme(topic as string, theme as string, token)

            if (data?.error) {
                console.error(data.error)
                alert(data.error)
                router.refresh()
                setGlobalLoading(false)
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
    }, [topic, theme, token, router])

    const checkAnswer = async () => {
        if (checkLoading || tasks[activeTask].isDone || code == "") {
            return
        }

        setCheckLoading(true)

        const data = await resolve(tasks[activeTask].id, JSON.stringify(code).slice(1, -1).replace(/\n/g, '\\n'), token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            setCheckLoading(false)
            return
        }

        setTasks(ts => {
            if (!ts[activeTask].dialog.includes(data.message)) {
                ts[activeTask].dialog.push(data.message)
            }

            ts[activeTask].isDone = data.ok
            return ts
        })

        setCheckLoading(false)
    }

    const compilation = async () => {
        if (compilationLoading || code == "") {
            return
        }

        setCompilationLoading(true)

        const data = await compile(JSON.stringify(code).slice(1, -1).replace(/\n/g, '\\n'), token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            setCompilationLoading(false)

            return
        }

        if (data?.result) {
            setConsoleText(p => [...p, data.result])
        }

        setCompilationLoading(false)
    }

    const loadMoreProblems = async () => {
        if (problemsLoading) {
            return
        }

        setProblemsLoading(true)

        const data = await getProblems(topic as string, theme as string, token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            setProblemsLoading(false)
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


        setProblemsLoading(false)
    }

    const canNextTask = useMemo(() => tasks.length > activeTask + 1, [tasks.length, activeTask])

    const nextTask = () => {
        if (!canNextTask || checkLoading) {
            return
        }

        setActiveTask(p => p + 1)
    }

    const canPrevTask = useMemo(() => activeTask > 0, [tasks.length, activeTask])

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
                        <button onClick={() => compilation()} disabled={compilationLoading || code == ""} type="button" className="btn btn-success btn-sm position-absolute" style={{ top: '10px', right: '10px', zIndex: 999 }}>{"▶"}</button>
                        <div className="border" style={{ height: 'calc(100% - 150px)', overflow: 'scroll' }}>
                            <CodeMirror
                                onChange={setCode}
                                value={code}
                                height="100%"
                            />
                        </div>
                        <div className="border p-2 " style={{ height: '150px', overflow: 'scroll' }}>
                            {consoleText.map((s, i) => <p key={i}>{s}</p>)}
                            {compilationLoading && <p>Думаю что это...</p>}
                        </div>
                    </div>
                    <div className="col-6 border rounded p-2" style={{ height: '100%', overflow: 'scroll' }}>
                        {part === 'tasks'
                            ? <>
                                <div className="text p-2 border rounded" style={{ height: 'calc(100% - 50px)', overflow: 'scroll' }}>
                                    {tasks.length > 0 && <div dangerouslySetInnerHTML={{ __html: tasks[activeTask].task }} />}
                                    {tasks.length > 0 && tasks[activeTask].dialog.map((d, i) => <p key={i}><strong>AI:</strong> {d}</p>)}
                                    {checkLoading && <p><strong>AI:</strong> Думаю...</p>}
                                </div>
                                <div className="p-2 border rounded d-flex gap-1 align-items-center" style={{ height: '50px' }}>
                                    {tasks.length > 0 && <button onClick={() => prevTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canPrevTask}>{"<<"}</button>}
                                    {tasks.length > 0 && <button disabled={checkLoading || tasks[activeTask].isDone || code == ""} onClick={() => checkAnswer()} type="button" className="btn btn-outline-success btn-sm">Проверить решение</button>}
                                    {tasks.length > 0 && <button onClick={() => nextTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canNextTask}>{">>"}</button>}
                                    {tasks.length > 0 && <p className="test-center mx-auto mb-0">{activeTask + 1} / {tasks.length}</p>}
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
