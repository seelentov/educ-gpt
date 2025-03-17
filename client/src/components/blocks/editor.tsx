'use client'

import { getProblems } from "@/core/api/roadmap/problems"
import { getTheme } from "@/core/api/roadmap/theme"
import { useLocalStorage } from "@/core/hooks/useLocalStorage"
import { useParams, useRouter } from "next/navigation"
import { useState, useEffect, useMemo } from "react"
import { Loading } from "../ui/loading"
import { remark } from 'remark';
import html from 'remark-html';
import rehypeHihglight from 'rehype-highlight'
import { resolve } from "@/core/api/roadmap/resolve"
import { compile } from "@/core/api/utils/compile"
import CodeMirror from "@uiw/react-codemirror";
import { checkAnswerUtil } from "@/core/api/utils/check_answer"
import Select from 'react-select'

interface Task {
    task: string, isDone: boolean, dialog: string[], id: number, isTheory: boolean, languages: string[]
}

export function Editor() {
    const [part, setPath] = useState<"theory" | "tasks">("theory")
    const [content, setContent] = useState<string>("")
    const [tasks, setTasks] = useState<Task[]>([])
    const [activeTask, setActiveTask] = useState<number>(0)
    const [code, setCode] = useState<string>("")
    const [activeLanguage, setActiveLanguage] = useState<string>("")
    const [consoleText, setConsoleText] = useState<string[]>([])
    const { topic, theme } = useParams()
    const [globalLoading, setGlobalLoading] = useState<boolean>(false)
    const [problemsLoading, setProblemsLoading] = useState<boolean>(false)
    const [checkLoading, setCheckLoading] = useState<boolean>(false)
    const [compilationLoading, setCompilationLoading] = useState<boolean>(false)

    const [token] = useLocalStorage("token", "")

    const router = useRouter()

    useEffect(() => {
        if (token === "") {
            router.replace("/login")
        }
    })

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
                    isTheory: p.is_theory,
                    languages: p.languages.split(";"),
                    dialog: []
                }
            }))
            setGlobalLoading(false)
        })()
    }, [topic, theme, token, router])

    const languageOptions = useMemo(() => {
        if (part == 'theory') {
            return "Python; JavaScript; Java; C++; C#; PHP; TypeScript; Swift; Go; Kotlin; Ruby; Rust; SQL; R; Perl; Dart; Scala; Haskell; Lua; Objective-C; Shell; PowerShell; Assembly; MATLAB; Groovy; Elixir; Clojure; F#; Erlang; VBA; Delphi; Ada; Lisp; Fortran; Prolog; Cobol; Bash; Racket; Julia; Crystal; Nim; OCaml; D; Vala; Smalltalk; ABAP; ActionScript; Apex; ColdFusion; Eiffel; LabVIEW; PL/SQL; SAS; Scheme; Tcl; Verilog; VHDL; Zig".split("; ").map((l) => ({
                value: l,
                label: l
            }))
        }


        return tasks[activeTask]?.languages.map(l => ({
            value: l,
            label: l
        })) || []
    }, [tasks, activeTask, part])

    useEffect(() => {
        if (languageOptions.length > 0) {
            setActiveLanguage(languageOptions[0].value)
        }
    }, [languageOptions])

    useEffect(() => {
        if (part === 'tasks' && tasks[activeTask].isTheory) {
            setActiveLanguage("")
        }
    }, [part])

    const checkAnswer = async () => {
        if (checkLoading || code == "") {
            return
        }

        setCheckLoading(true)

        let data

        if (tasks[activeTask].isDone) {
            data = await checkAnswerUtil(tasks[activeTask].task.replace(/<[^>]*>?/gm, '').replace(/\n/g, ' '), JSON.stringify(code).slice(1, -1).replace(/\n/g, '\\n'), activeLanguage || "", token)
        } else {
            data = await resolve(tasks[activeTask].id, JSON.stringify(code).slice(1, -1).replace(/\n/g, '\\n'), activeLanguage || "", token)
        }

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

            if (!ts[activeTask].isDone) {
                ts[activeTask].isDone = data.ok
            }

            return ts
        })

        setCheckLoading(false)
    }

    const compilation = async () => {
        if (compilationLoading || code == "") {
            return
        }

        setCompilationLoading(true)

        const data = await compile(JSON.stringify(code).slice(1, -1).replace(/\n/g, '\\n'), activeLanguage || "", token)

        if (data?.error) {
            console.error(data.error)
            alert(JSON.stringify(data.error))
            setCompilationLoading(false)

            return
        }

        if (data?.result) {
            setConsoleText(p => [...p, data.result.split("\n")])
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
                    id: p.id,
                    task: p.question,
                    isDone: false,
                    isTheory: p.is_theory,
                    languages: p.languages.split(";"),
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

    const leftSizeIsTheory = useMemo(() => (tasks && tasks.length > 0 && tasks[activeTask].isTheory && part === 'tasks'), [tasks, activeTask, part])


    return (
        <>
            {!globalLoading
                ? <div className="editor-container col12 p-3 border d-flex gap-3 pt-5 position-relative" style={{ height: '100vh' }}>
                    <div className="position-absolute fixed-top p-2 d-flex" style={{ height: '48px' }}>
                        <button onClick={() => router.back()} type="button" className={`btn btn-outline-primary btn-sm me-2`}>Назад</button>
                        <button onClick={() => setPath('theory')} style={{ marginLeft: 'auto' }} type="button" className={`btn ${part === 'theory' ? "btn-primary" : "btn-outline-primary"} btn-sm me-2`}>Теория</button>
                        <button onClick={() => setPath('tasks')} type="button" className={`btn ${part === 'tasks' ? "btn-primary" : "btn-outline-primary"} btn-sm me-2`}>Задачи</button>
                    </div>
                    <div className="editor-left col-6 border rounded position-relative" style={{ height: '100%' }}>
                        {!leftSizeIsTheory && <div className="border d-flex gap-1 align-items-flex-start" style={{ height: '50px', padding: '5px' }}>
                            <button onClick={() => compilation()} disabled={compilationLoading || code == ""} type="button" className="btn btn-success" style={{ zIndex: 999 }}>{"▶"}</button>
                            {languageOptions.length > 0 && (
                                <Select
                                    className="basic-single"
                                    classNamePrefix="select"
                                    value={languageOptions.find(option => option.value === activeLanguage)}
                                    onChange={(selectedOption) => setActiveLanguage(selectedOption?.value || "")}
                                    options={languageOptions}
                                    isSearchable={true}
                                    styles={{
                                        container: (provided) => ({
                                            ...provided,
                                            width: '150px',
                                            height: '30px',
                                            fontSize: '13px'
                                        }),
                                        control: (provided) => ({
                                            ...provided,
                                            width: '100%',
                                        }),
                                        menu: (provided) => ({
                                            ...provided,
                                            width: '150px',
                                        }),
                                    }}
                                />
                            )}
                        </div>}
                        <div className={leftSizeIsTheory ? "" : "border"} style={{ height: leftSizeIsTheory ? '100%' : 'calc(100% - 190px)', overflowY: leftSizeIsTheory ? "initial" : 'scroll' }}>
                            {
                                leftSizeIsTheory
                                    ? <textarea
                                        onChange={(e) => setCode(e.target.value)}
                                        value={code}
                                        className="p-3 w-100 h-100"
                                        style={{ resize: 'none', overflowY: 'scroll', height: '100%' }}
                                        placeholder="Ваш ответ..."
                                    />
                                    : <CodeMirror
                                        onChange={setCode}
                                        value={code}
                                        height="1000px"
                                    />
                            }
                        </div>
                        {!leftSizeIsTheory
                            &&
                            <div className="console-container border p-2 " style={{ overflowY: 'scroll' }}>
                                {consoleText.map((s, i) => <p key={i}>{s}</p>)}
                                {compilationLoading && <p>Думаю что это...</p>}
                            </div>}
                    </div>
                    <div className="editor-right col-6 border rounded p-2" style={{ height: '100%', overflowY: 'scroll' }}>
                        {part === 'tasks'
                            ? <>
                                <div className="text p-2 border rounded" style={{ height: 'calc(100% - 50px)', overflowY: 'scroll' }}>
                                    {tasks[activeTask].isDone && <p className="text-success">Ответ принят</p>}
                                    {tasks.length > 0 && <div dangerouslySetInnerHTML={{ __html: tasks[activeTask].task }} />}
                                    {tasks.length > 0 && tasks[activeTask].dialog.map((d, i) => <p key={i}><strong>AI:</strong> {d}</p>)}
                                    {checkLoading && <p><strong>AI:</strong> Думаю...</p>}
                                </div>
                                <div className="p-2 border rounded d-flex flex-column gap-1 align-items-center editor-btns">
                                    <div className="d-flex gap-1">
                                        {tasks.length > 0 && <button onClick={() => prevTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canPrevTask}>{"<<"}</button>}
                                        {tasks.length > 0 && <button onClick={() => nextTask()} type="button" className="btn btn-outline-primary btn-sm" disabled={!canNextTask}>{">>"}</button>}
                                    </div>
                                    {tasks.length > 0 && <button disabled={checkLoading || code == ""} onClick={() => checkAnswer()} type="button" className="btn btn-outline-success btn-sm w-100">{tasks[activeTask].isDone ? "Внести правки" : "Проверить решение"}</button>}
                                    {tasks.length > 0 && <p className="text-center mx-auto mb-0 text-nowrap">{activeTask + 1} / {tasks.length}</p>}
                                    <button onClick={() => loadMoreProblems()} disabled={problemsLoading} type="button" className="btn btn-outline-warning btn-sm w-100">Загрузить еще...</button>
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