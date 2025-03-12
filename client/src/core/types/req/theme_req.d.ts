interface ThemeRequest {
    text: string
    problems: Problem[]
}

interface Problem {
    id: number
    question: string
}