import { API_URL } from "./api";

export async function baseFetch(
    url: string,
    body: string | FormData | null = null,
    method: string = "GET",
    token: string = "",
    headers: { [key: string]: string } | null | undefined = null,
    no_json: boolean = false
): Promise<any> {
    try {
        const init: RequestInit = {
            method,
            headers: {}
        };

        if (!no_json) {
            init.headers = {
                ...init.headers,
                'Content-Type': 'application/json',
            }
        }

        if (headers) {
            init.headers = {
                ...init.headers,
                ...headers,
            };
        }

        if (token !== "") {
            init.headers = {
                ...init.headers,
                "Authorization": "Bearer " + token
            }
        }

        if (body) {
            init.body = body;
        }

        const response = await fetch(`${API_URL}${url}`, init);

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Ошибка:', error);

        if (error instanceof Error) {
            throw error;
        }

        throw new Error(`Неизвестная ошибка: ${error}`);
    }
}