

import { useState } from "react"

export const useLocalStorage = (key: string, initialValue: string) => {
    const [state, setState] = useState(() => {
        try {
            const value = window.localStorage.getItem(key)
            return value ? JSON.parse(value) : initialValue
        } catch (error) {
            console.log(error)
        }
    })

    const setValue = (value: any) => {
        try {
            const valueToStore = value instanceof Function ? value(state) : value
            window.localStorage.setItem(key, JSON.stringify(valueToStore))
            setState(value)
        } catch (error) {
            console.log(error)
        }
    }

    return [state, setValue]
}