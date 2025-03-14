interface ThemeRequest {
    text: string
    problems: Problem[]
}

interface Problem {
    id: number
    question: string
    languages: string
    is_theory: boolean
}